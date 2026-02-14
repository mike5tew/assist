#!/usr/bin/env python3
"""Add 'Reception - Sticker Book' course to Weaviate CourseSkillSuggestions."""

import json
import urllib.request

WEAVIATE_URL = "http://localhost:8081"

course_obj = {
    "class": "CourseSkillSuggestions",
    "properties": {
        "course": "Reception - Sticker Book",
        "course_code": "REC-SB",
        "subject": "Primary",
        "key_stage": "EYFS",
        "year_group": "Reception",
        "description": "A sticker-book course for Reception that maps foundational developmental skills across cognitive, motor, social and communication domains.",
        "core_skills": [
            "Attention",
            "Working Memory",
            "Fine Motor Skills",
            "Following Protocols",
            "Social Confidence",
            "Speaking",
            "Vocabulary",
            "Non-verbal communication",
            "Self-regulation",
            "Empathy",
        ],
        "practice_skills": [],
        "implicit_skills": [],
        "all_suggested_skills": "Attention, Working Memory, Fine Motor Skills, Following Protocols, Social Confidence, Speaking, Vocabulary, Non-verbal communication, Self-regulation, Empathy",
    },
}

def main():
    # Check if it already exists
    check_query = json.dumps({"query": '{ Get { CourseSkillSuggestions(where: { path: ["course"], operator: Equal, valueText: "Reception - Sticker Book" }, limit: 1) { course subject _additional { id } } } }'}).encode()
    req = urllib.request.Request(f"{WEAVIATE_URL}/v1/graphql", data=check_query, headers={"Content-Type": "application/json"})
    with urllib.request.urlopen(req) as resp:
        existing = json.loads(resp.read())

    existing_courses = existing["data"]["Get"]["CourseSkillSuggestions"]
    if existing_courses:
        print(f"Already exists with ID: {existing_courses[0]['_additional']['id']}")
        print("Skipping creation.")
        return

    # Create the object
    payload = json.dumps(course_obj).encode()
    req = urllib.request.Request(f"{WEAVIATE_URL}/v1/objects", data=payload, headers={"Content-Type": "application/json"})
    with urllib.request.urlopen(req) as resp:
        result = json.loads(resp.read())

    print(f"Created successfully!")
    print(f"  ID:         {result.get('id')}")
    print(f"  Class:      {result.get('class')}")
    props = result.get("properties", {})
    print(f"  Course:     {props.get('course')}")
    print(f"  Subject:    {props.get('subject')}")
    print(f"  Key Stage:  {props.get('key_stage')}")
    print(f"  Year Group: {props.get('year_group')}")
    print(f"  Core Skills ({len(props.get('core_skills', []))}): {props.get('core_skills')}")

if __name__ == "__main__":
    main()
