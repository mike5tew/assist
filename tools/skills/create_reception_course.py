#!/usr/bin/env python3
"""
create_reception_course.py

Safe, idempotent script to create (if missing) the Subject "Primary" and Course "Reception - Sticker Book",
and attach existing skills referenced in `assist/docs/training_data/RECEPTION_STICKER_BOOK.json` to that Course.

Usage:
  # Dry-run (no DB writes):
  python3 create_reception_course.py --dry-run

  # Apply changes (create subject/course and attach skills):
  python3 create_reception_course.py --apply

The script will NOT create new skills. Any skills not found in `skills_key` are listed in the report for manual integration.

DB connectivity: by default the script reads DB connection parts from environment variables:
  DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME
Alternatively set DATABASE_URL in the form: user:password@tcp(host:port)/dbname

"""

import argparse
import json
import os
import re
import sys
from typing import List, Tuple

try:
    import pymysql
except Exception as e:
    print("Missing dependency: pymysql. Install with: pip install pymysql")
    raise

# Defaults
DEFAULT_SUBJECT_NAME = "Primary"
DEFAULT_COURSE_NAME = "Reception - Sticker Book"
JSON_PATH = os.path.join(os.path.dirname(__file__), '..', '..', 'docs', 'training_data', 'RECEPTION_STICKER_BOOK.json')


def parse_database_url(url: str):
    # support quick parse of user:pass@tcp(host:port)/dbname?...
    m = re.match(r"(?P<user>[^:]+):(?P<pw>[^@]+)@tcp\((?P<host>[^:]+):(?P<port>\d+)\)/(?P<db>[^?]+)", url)
    if not m:
        raise ValueError("Unsupported DATABASE_URL format. Use user:password@tcp(host:port)/dbname")
    return m.group('user'), m.group('pw'), m.group('host'), int(m.group('port')), m.group('db')


def get_db_connection():
    # Allow DATABASE_URL or individual env vars
    DATABASE_URL = os.getenv('DATABASE_URL')
    if DATABASE_URL:
        user, pw, host, port, db = parse_database_url(DATABASE_URL)
    else:
        user = os.getenv('DB_USER', 'root')
        pw = os.getenv('DB_PASSWORD', 'skills_password123')
        host = os.getenv('DB_HOST', '127.0.0.1')
        port = int(os.getenv('DB_PORT', '3306'))
        db = os.getenv('DB_NAME', 'dare2lead')

    conn = pymysql.connect(host=host, port=port, user=user, password=pw, database=db, autocommit=True)
    return conn


def load_json(path: str):
    with open(path, 'r', encoding='utf-8') as f:
        return json.load(f)


def find_subject(conn, name: str) -> int:
    with conn.cursor() as cur:
        cur.execute("SELECT SubjectID FROM subject WHERE LOWER(Subject_name)=LOWER(%s) LIMIT 1", (name,))
        row = cur.fetchone()
        return int(row[0]) if row else 0


def create_subject(conn, name: str, code: str = '') -> int:
    with conn.cursor() as cur:
        cur.execute("INSERT INTO subject (Subject_name, Subject_code) VALUES (%s, %s)", (name, code))
        return cur.lastrowid


def find_course(conn, name: str, subject_id: int) -> int:
    with conn.cursor() as cur:
        cur.execute("SELECT CourseID FROM courses WHERE LOWER(Course_name)=LOWER(%s) AND SubjectID=%s LIMIT 1", (name, subject_id))
        row = cur.fetchone()
        return int(row[0]) if row else 0


def create_course(conn, name: str, subject_id: int) -> int:
    with conn.cursor() as cur:
        cur.execute("INSERT INTO courses (Course_name, Course_code, SubjectID) VALUES (%s, %s, %s)", (name, '', subject_id))
        return cur.lastrowid


def find_skill_by_name(conn, name: str) -> List[Tuple[int, str]]:
    with conn.cursor() as cur:
        cur.execute("SELECT Skills_keyID, Skill_name FROM skills_key WHERE LOWER(Skill_name)=LOWER(%s)", (name,))
        rows = cur.fetchall()
        return [(int(r[0]), r[1]) for r in rows]


def create_skill_attached(conn, skill_id: int, course_id: int):
    with conn.cursor() as cur:
        cur.execute("INSERT INTO skill_attached (Skills_keyID, Lesson_FilesID, CourseID) VALUES (%s, %s, %s)", (skill_id, 0, course_id))
        return cur.lastrowid


def check_skill_attached_exists(conn, skill_id: int, course_id: int) -> bool:
    with conn.cursor() as cur:
        cur.execute("SELECT Skill_attachedID FROM skill_attached WHERE Skills_keyID=%s AND CourseID=%s LIMIT 1", (skill_id, course_id))
        return cur.fetchone() is not None


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--dry-run', action='store_true', default=True, dest='dry_run', help='Dry run (no DB writes). Default: on')
    parser.add_argument('--apply', action='store_true', dest='apply', help='Apply changes (create subject/course and attach skills)')
    parser.add_argument('--subject-name', default=DEFAULT_SUBJECT_NAME)
    parser.add_argument('--course-name', default=DEFAULT_COURSE_NAME)
    parser.add_argument('--json-path', default=JSON_PATH)

    args = parser.parse_args()

    if not args.apply:
        print("Running in dry-run mode. No changes will be made. Use --apply to persist.")

    data = load_json(args.json_path)
    skills = [s.get('name') for s in data.get('skills', [])]

    conn = get_db_connection()

    try:
        # Subject
        subject_id = find_subject(conn, args.subject_name)
        if subject_id:
            print(f"Found subject '{args.subject_name}' (SubjectID={subject_id})")
        else:
            print(f"Subject '{args.subject_name}' not found")
            if args.apply:
                subject_id = create_subject(conn, args.subject_name)
                print(f"Created subject '{args.subject_name}' (SubjectID={subject_id})")
            else:
                print("(dry-run) would create subject")

        if not subject_id:
            print("ERROR: subject missing and not created. Use --apply to create.")
            sys.exit(1)

        # Course
        course_id = find_course(conn, args.course_name, subject_id)
        if course_id:
            print(f"Found course '{args.course_name}' (CourseID={course_id})")
        else:
            print(f"Course '{args.course_name}' not found under subject '{args.subject_name}'")
            if args.apply:
                course_id = create_course(conn, args.course_name, subject_id)
                print(f"Created course '{args.course_name}' (CourseID={course_id})")
            else:
                print("(dry-run) would create course")

        if not course_id:
            print("ERROR: course missing and not created. Use --apply to create.")
            sys.exit(1)

        # Attach skills
        attached = []
        missing = []
        skipped = []
        for skill_name in skills:
            found = find_skill_by_name(conn, skill_name)
            if not found:
                missing.append(skill_name)
                print(f"MISSING: skill not found in skills_key -> '{skill_name}'")
                continue
            # if multiple matches pick the exact case-insensitive one first
            skill_id = found[0][0]
            skill_db_name = found[0][1]
            if check_skill_attached_exists(conn, skill_id, course_id):
                skipped.append((skill_id, skill_db_name))
                print(f"SKIP (already attached): {skill_db_name} (Skills_keyID={skill_id})")
                continue
            if args.apply:
                sid = create_skill_attached(conn, skill_id, course_id)
                attached.append((skill_id, skill_db_name, sid))
                print(f"ATTACHED: {skill_db_name} (Skills_keyID={skill_id}) as Skill_attachedID={sid}")
            else:
                print(f"(dry-run) would attach: {skill_db_name} (Skills_keyID={skill_id}) to course {course_id}")

        # Summary
        print('\nSUMMARY:')
        print(f" Total skills processed: {len(skills)}")
        print(f" Attached: {len(attached)}")
        print(f" Skipped (already attached): {len(skipped)}")
        print(f" Missing (not in skills_key): {len(missing)}")
        if missing:
            print("Missing skills list:")
            for m in missing:
                print(' -', m)

        if missing and not args.apply:
            print('\nNote: Some skills above were not found in skills_key. You can either add them to skills_key (so they are real platform skills) or choose alternatives. The script intentionally does NOT create new skills.')

    finally:
        conn.close()


if __name__ == '__main__':
    main()
