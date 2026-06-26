#!/usr/bin/env python3
"""
ToddlerOS Book 1 -- "Hello World" PDF Generator

Generates the first edition activity book: all 12 spectra at Cycle 1.
Content sourced from TODDLEROS_BOOKLET_OUTLINES.md

Landscape A5, warm cream palette, Trebuchet MS body, Georgia headings.
"""

from fpdf import FPDF
import os
import re

# Design palette
CREAM = (255, 248, 240)
TAN = (245, 230, 211)
DARK_BROWN = (62, 39, 35)
WARM_BEIGE = (232, 213, 183)
LEAF_GREEN = (139, 195, 74)
SOFT_GREY = (180, 170, 160)
WHITE = (255, 255, 255)
GREEN_BG = (220, 245, 220)
RED_BG = (245, 220, 220)

OUTPUT_DIR = os.path.dirname(os.path.abspath(__file__))
OUTPUT_PATH = os.path.join(OUTPUT_DIR, "..", "docs", "ToddlerOS_Book1_HelloWorld.pdf")

SYSTEM_FONT_DIR = "/System/Library/Fonts/Supplemental/"
ASSET_DIR = "/Users/michaelstewart/Coding/ToddlerOS/assets/"
BLAPEROOM = ASSET_DIR + "Blaperoom/OpenType-TT/BLAPEROOM.ttf"
GILL_SANS = ASSET_DIR + "gill-sans-infant-std/gill-sans-infant-std.otf"


def strip_emoji(text):
    """Remove emoji but keep em-dashes and standard punctuation."""
    emoji_pattern = re.compile(
        "[\U0001F300-\U0001F9FF"
        "\U0001FA00-\U0001FAFF"
        "\u2600-\u26FF"
        "\u2700-\u27BF"
        "\uFE00-\uFE0F"
        "\u200D"
        "]+",
        flags=re.UNICODE,
    )
    return emoji_pattern.sub("", text).strip()


class ToddlerOSBook(FPDF):
    def __init__(self):
        super().__init__(orientation="L", unit="mm", format=(148, 210))
        self.set_auto_page_break(auto=False)
        self.set_margins(12, 12, 12)
        # Body: Gill Sans Infant Std (regular), Trebuchet MS (bold/italic fallback)
        self.add_font("body", "", GILL_SANS, uni=True)
        self.add_font("body", "B", SYSTEM_FONT_DIR + "Trebuchet MS Bold.ttf", uni=True)
        self.add_font("body", "I", SYSTEM_FONT_DIR + "Trebuchet MS Italic.ttf", uni=True)
        self.add_font("body", "BI", SYSTEM_FONT_DIR + "Trebuchet MS Bold Italic.ttf", uni=True)
        # Heading: BLAPEROOM display font
        self.add_font("heading", "", BLAPEROOM, uni=True)
        self.add_font("heading", "B", BLAPEROOM, uni=True)
        self.add_font("heading", "I", BLAPEROOM, uni=True)

    def _bg(self, color=CREAM):
        self.set_fill_color(*color)
        self.rect(0, 0, self.w, self.h, "F")

    def _header_bar(self, text, y=10):
        self.set_fill_color(*TAN)
        self.rect(0, y - 2, self.w, 12, "F")
        self.set_font("heading", "B", 13)
        self.set_text_color(*DARK_BROWN)
        self.set_xy(12, y)
        self.cell(self.w - 24, 8, strip_emoji(text), align="C")

    def _divider(self, y):
        self.set_draw_color(*WARM_BEIGE)
        self.set_line_width(0.3)
        self.line(12, y, self.w - 12, y)

    # -------- PAGES ------------------------------------------------

    def cover_page(self):
        self.add_page()
        self._bg(TAN)

        self.set_font("heading", "B", 28)
        self.set_text_color(*DARK_BROWN)
        self.set_xy(0, 25)
        self.cell(self.w, 12, "ToddlerOS", align="C")

        self.set_font("heading", "I", 18)
        self.set_xy(0, 42)
        self.cell(self.w, 8, "Hello World", align="C")

        self.set_font("body", "", 10)
        self.set_xy(0, 56)
        self.cell(self.w, 6, "Edition 1  --  All 12 Themes at Foundation Level", align="C")

        self.set_font("body", "", 8)
        self.set_xy(0, 70)
        self.cell(self.w, 5, "12 activity spreads  /  48 stickers  /  Parent reference card", align="C")

        themes = [
            "Near and Far", "Loud and Quiet", "Big Feelings, Small Feelings",
            "Yes and No", "Mine and Yours", "Try and Wait",
            "Same and Different", "Your Feelings, My Feelings", "Tidy and Messy",
            "My Fault, Your Fault", "Keeping and Letting Go", "Wanting and Waiting",
        ]
        y = 82
        for i in range(0, 12, 3):
            row = "   /   ".join(themes[i : i + 3])
            self.set_xy(0, y)
            self.cell(self.w, 4, row, align="C")
            y += 5

        self.set_font("body", "I", 7)
        self.set_xy(0, self.h - 14)
        self.cell(self.w, 4, "DRAFT -- For review purposes only -- April 2026", align="C")

    def the_dance_page(self):
        self.add_page()
        self._bg()

        self.set_font("heading", "B", 13)
        self.set_text_color(*DARK_BROWN)
        self.set_xy(14, 8)
        self.cell(self.w - 28, 7, "Before you start")

        blocks = [
            ("", "Your child has a way of being in the world. So do you."),
            ("", "Sometimes it's easy. Sometimes it's not. Some days they're brave and curious and kind. Some days they scream in the cereal aisle. Some days you handle it perfectly. Some days you don't."),
            ("", "This book isn't about fixing any of that."),
            ("", "There's a musician called Ren who said something that stuck with us:"),
            ("I", '"It was never a battle I was supposed to win. It is an eternal dance, and just like any dance, the more rigid I became, the harder it got. The more I cursed my clumsy footsteps, the more I struggled. So I got older, and I learned to relax, and soften, and that dance got easier."'),
            ("B", "That's what this book is about."),
            ("", "Your child has settings -- biological ones. How big their feelings are. Whether people charge their battery or drain it. Whether they jump in or watch first. These aren't problems. They're just settings."),
            ("", 'Your child also has skills -- things they can do fluently, and things they\'re still building. Not compared to other children. Not measured against "normal." Just: what can they do, and what\'s next?'),
            ("", "The settings and the skills together make a pattern. That pattern is unique to your child. It's not broken. It's not behind. It's theirs."),
            ("", "The activities in here aren't tests. There's no pass or fail. They're invitations -- little experiments where you and your child play together, and both of you notice something."),
            ("B", "The only rule: if it's a bad day, close the book. Read them a story instead. The dance will be there tomorrow."),
        ]

        y = 17
        for style, text in blocks:
            self.set_font("body", style, 8)
            self.set_xy(14, y)
            self.multi_cell(self.w - 28, 3.8, text, align="L")
            y = self.get_y() + 1.2

    def how_to_use_page(self):
        self.add_page()
        self._bg()
        self._header_bar("How This Book Works")

        sections = [
            ("Two modes per spread",
             'Each theme has two activities side by side.\n\n'
             '"With You" (2-5 min): You and your child play together. You\'re watching, naming, gently stretching.\n\n'
             '"While You..." (5-15 min): Your child works independently -- colouring, sticking, drawing.'),
            ("Three pathways",
             "Every activity ends with three possible outcomes:\n\n"
             ">> ESCALATE -- They showed range. Step forward.\n"
             "<< SIMPLIFY -- Not ready. Step back to something easier.\n"
             "!! VOLTAGE SPIKE -- Not a skill gap. Their emotional conductor overloaded. Stop. Name the feeling. Wait."),
            ("Reading order",
             "There isn't one. Open to any page. Every theme stands alone. Do them in any order."),
            ("Bad Day Rule",
             "If it's a bad day, close the book. Read them a story instead. Bad days are data about sleep, hunger, or timing -- not about your child."),
        ]

        y = 26
        for heading, body in sections:
            self.set_font("heading", "B", 9)
            self.set_text_color(*DARK_BROWN)
            self.set_xy(14, y)
            self.cell(self.w - 28, 5, heading)
            y += 5.5
            self.set_font("body", "", 7.5)
            self.set_xy(14, y)
            self.multi_cell(self.w - 28, 3.5, body, align="L")
            y = self.get_y() + 3

    def activity_spread(self, theme_num, theme_name, spectrum, observation,
                        cycle_label, activity, desc, duration,
                        notice, rescue, try_lang, avoid_lang,
                        escalate, simplify, stickers,
                        while_prompt, while_desc):
        # ---------- LEFT PAGE: WITH YOU ----------
        self.add_page()
        self._bg()

        # Header bar
        self.set_fill_color(*TAN)
        self.rect(0, 0, self.w, 18, "F")
        self.set_font("heading", "B", 11)
        self.set_text_color(*DARK_BROWN)
        self.set_xy(12, 3)
        self.cell(self.w - 70, 6, f"Theme {theme_num}: {theme_name}")
        self.set_font("body", "", 7.5)
        self.set_xy(12, 10)
        self.cell(self.w - 70, 5, f"Spectrum: {spectrum}  |  Cycle 1: {cycle_label}")
        self.set_font("body", "B", 8)
        self.set_xy(self.w - 55, 4)
        self.cell(43, 5, duration, align="R")
        self.set_font("body", "", 7)
        self.set_xy(self.w - 55, 10)
        self.cell(43, 5, "WITH YOU", align="R")

        y = 21

        # Activity name + description
        self.set_font("heading", "B", 10)
        self.set_xy(12, y)
        self.cell(self.w - 24, 5, strip_emoji(activity))
        y += 6
        self.set_font("body", "", 7.5)
        self.set_xy(12, y)
        self.multi_cell(self.w - 24, 3.5, strip_emoji(desc))
        y = self.get_y() + 2
        self._divider(y); y += 2

        # What to notice
        self.set_font("body", "B", 7.5)
        self.set_xy(12, y)
        self.cell(self.w - 24, 4, "WHAT TO NOTICE")
        y += 4.5
        self.set_font("body", "", 7)
        self.set_xy(12, y)
        self.multi_cell(self.w - 24, 3.3, strip_emoji(notice))
        y = self.get_y() + 1.5

        # Rescue
        self.set_font("body", "B", 7.5)
        self.set_xy(12, y)
        self.cell(self.w - 24, 4, "IF THEY STOP / STRUGGLE...")
        y += 4.5
        self.set_font("body", "", 7)
        self.set_xy(12, y)
        self.multi_cell(self.w - 24, 3.3, strip_emoji(rescue))
        y = self.get_y() + 1.5

        # Try / Avoid language boxes
        if y < self.h - 38:
            col_w = (self.w - 28) / 2
            box_h = 10

            self.set_fill_color(*GREEN_BG)
            self.rect(12, y, col_w, box_h, "F")
            self.set_font("body", "B", 6.5)
            self.set_xy(13, y + 0.5)
            self.cell(col_w - 2, 3, "TRY:")
            self.set_font("body", "I", 6.5)
            self.set_xy(13, y + 4)
            self.multi_cell(col_w - 4, 2.8, f'"{strip_emoji(try_lang)}"')

            x2 = 12 + col_w + 4
            self.set_fill_color(*RED_BG)
            self.rect(x2, y, col_w, box_h, "F")
            self.set_font("body", "B", 6.5)
            self.set_xy(x2 + 1, y + 0.5)
            self.cell(col_w - 2, 3, "NOT:")
            self.set_font("body", "I", 6.5)
            self.set_xy(x2 + 1, y + 4)
            self.multi_cell(col_w - 4, 2.8, f'"{strip_emoji(avoid_lang)}"')
            y += box_h + 2

        # Pathways
        if y < self.h - 22:
            self.set_font("body", "B", 7)
            self.set_xy(12, y)
            self.cell(self.w - 24, 4, "NEXT STEP")
            y += 4.5
            self.set_font("body", "", 6.5)
            for arrow, label, txt in [
                (">>", "ESCALATE", strip_emoji(escalate)),
                ("<<", "SIMPLIFY", strip_emoji(simplify)),
                ("!!", "VOLTAGE SPIKE", "Stop. Name the feeling. Wait for the wave."),
            ]:
                self.set_xy(14, y)
                self.set_font("body", "B", 6.5)
                self.cell(5, 3.2, arrow)
                self.cell(22, 3.2, label)
                self.set_font("body", "", 6.5)
                self.cell(self.w - 60, 3.2, f" -- {txt}")
                y += 4

        # Footer observation prompt
        self.set_font("body", "I", 6.5)
        self.set_text_color(*SOFT_GREY)
        self.set_xy(12, self.h - 9)
        self.cell(self.w - 24, 4, f'Observation: "{strip_emoji(observation)}"', align="C")

        # ---------- RIGHT PAGE: WHILE YOU ----------
        self.add_page()
        self._bg()

        self.set_fill_color(*TAN)
        self.rect(0, 0, self.w, 14, "F")
        self.set_font("heading", "B", 11)
        self.set_text_color(*DARK_BROWN)
        self.set_xy(12, 3)
        self.cell(self.w - 80, 6, f"While You...  --  {theme_name}")
        self.set_font("body", "", 7)
        self.set_xy(self.w - 72, 3)
        self.cell(60, 6, "5-15 min  |  CHILD INDEPENDENT", align="R")

        y = 18
        self.set_font("heading", "B", 10)
        self.set_xy(12, y)
        self.cell(self.w - 24, 5, strip_emoji(while_prompt))
        y += 7
        self.set_font("body", "", 7.5)
        self.set_xy(12, y)
        self.multi_cell(self.w - 24, 3.5, strip_emoji(while_desc))
        y = self.get_y() + 4

        # Illustration placeholder
        box_h = self.h - y - 30
        if box_h > 15:
            self.set_draw_color(*WARM_BEIGE)
            self.set_line_width(0.5)
            self.rect(20, y, self.w - 40, box_h, "D")
            self.set_font("body", "I", 9)
            self.set_text_color(*SOFT_GREY)
            self.set_xy(20, y + box_h / 2 - 4)
            self.cell(self.w - 40, 8, "[ Illustration space ]", align="C")
            y += box_h + 4

        # Stickers
        self.set_font("body", "B", 7.5)
        self.set_text_color(*DARK_BROWN)
        self.set_xy(12, y)
        self.cell(self.w - 24, 4, "This week, celebrate:")
        y += 4.5
        self.set_font("body", "", 7)
        for s in stickers:
            self.set_xy(16, y)
            self.cell(self.w - 32, 3.5, strip_emoji(s))
            y += 4

    def reference_card(self):
        self.add_page()
        self._bg(TAN)
        self._header_bar("Parent Reference Card")

        y = 26
        self.set_font("heading", "B", 9)
        self.set_text_color(*DARK_BROWN)
        self.set_xy(12, y)
        self.cell(self.w - 24, 5, "Language Shifts")
        y += 6

        shifts = [
            ('"You are..."', '"Right now you\'re..."'),
            ('"You should..."', '"Your brain is saying..."'),
            ('"That was wrong"', '"That didn\'t work. What happened?"'),
            ('"Try harder"', '"What\'s blocking you right now?"'),
            ('"Stop it"', '"I notice you\'re in X setting -- is that helping?"'),
            ('"Good boy/girl"', '"That worked. You figured it out."'),
            ('"Don\'t cry"', '"That\'s a big feeling. I\'m here."'),
        ]

        col_w = (self.w - 28) / 2
        self.set_fill_color(*WARM_BEIGE)
        self.rect(12, y, col_w, 5, "F")
        self.rect(14 + col_w, y, col_w, 5, "F")
        self.set_font("body", "B", 7)
        self.set_xy(13, y + 0.5)
        self.cell(col_w - 2, 4, "Old (Moral)")
        self.set_xy(15 + col_w, y + 0.5)
        self.cell(col_w - 2, 4, "New (Engineering)")
        y += 6

        self.set_font("body", "", 6.5)
        for old, new in shifts:
            self.set_xy(13, y)
            self.cell(col_w - 2, 3.5, old)
            self.set_xy(15 + col_w, y)
            self.cell(col_w - 2, 3.5, new)
            y += 4.5

        y += 2
        self._divider(y); y += 3

        self.set_font("heading", "B", 9)
        self.set_xy(12, y)
        self.cell(self.w - 24, 5, "Play-First Principle")
        y += 6
        steps = [
            '1. OBSERVE the setting without judgement',
            '2. NAME it neutrally ("You\'re in Independent mode right now")',
            '3. PLAY the activity that gently expands range',
            '4. NOTICE what happened ("When you had to wait, what happened?")',
        ]
        self.set_font("body", "", 7)
        for s in steps:
            self.set_xy(14, y)
            self.cell(self.w - 28, 3.5, s)
            y += 4.5

        y += 2
        self._divider(y); y += 3

        self.set_font("heading", "B", 9)
        self.set_xy(12, y)
        self.cell(self.w - 24, 5, "The Three Pathways")
        y += 6
        for label, txt in [
            (">> ESCALATE", "They showed range -- step forward on the skills graph."),
            ("<< SIMPLIFY", "Not ready -- drop to the foundation this needs."),
            ("!! VOLTAGE SPIKE", "Not a skill gap. Emotional overload. Stop. Ground. Wait."),
        ]:
            self.set_font("body", "B", 7)
            self.set_xy(14, y)
            self.cell(30, 3.5, label)
            self.set_font("body", "", 7)
            self.cell(self.w - 58, 3.5, txt)
            y += 5

    # ================================================================
    # CYOA: Jem's Big Day Out
    # ================================================================

    def _story_divider(self):
        """Section break page between activities and CYOA."""
        self.add_page()
        self._bg(TAN)
        self.set_font("heading", "B", 22)
        self.set_text_color(*DARK_BROWN)
        self.set_xy(0, 35)
        self.cell(self.w, 10, "Story Time", align="C")
        self.set_font("body", "I", 10)
        self.set_xy(0, 52)
        self.cell(self.w, 6, "A story where YOUR child chooses what happens", align="C")
        self._divider(64)
        self.set_font("body", "", 8)
        self.set_xy(20, 70)
        self.multi_cell(self.w - 40, 4,
            "The activities you've just done explore your child's settings.\n\n"
            "This story lets them USE those settings. Jem goes to a party. "
            "Your child decides what Jem does.\n\n"
            "There are no wrong answers. Every path leads somewhere.",
            align="C")

    def _cyoa_page(self, page_num, title, body_lines, choices, illustration_hint,
                   level_label=None):
        """Generic CYOA story page.

        body_lines: list of (style, text) tuples
        choices: list of (label, dest_page) tuples
        """
        self.add_page()
        self._bg()

        # Header
        self.set_fill_color(*TAN)
        self.rect(0, 0, self.w, 14, "F")
        self.set_font("heading", "B", 11)
        self.set_text_color(*DARK_BROWN)
        self.set_xy(12, 3)
        self.cell(self.w - 60, 6, title)
        self.set_font("body", "", 7)
        self.set_xy(self.w - 50, 3)
        right_label = f"Page {page_num}"
        if level_label:
            right_label += f"  |  {level_label}"
        self.cell(38, 6, right_label, align="R")

        y = 18

        # Body text
        for style, text in body_lines:
            self.set_font("body", style, 8)
            self.set_text_color(*DARK_BROWN)
            self.set_xy(14, y)
            self.multi_cell(self.w - 28, 3.8, strip_emoji(text))
            y = self.get_y() + 1.5

        # Illustration placeholder
        if y < self.h - 55:
            illo_h = min(25, self.h - y - 50)
            if illo_h > 10:
                self.set_draw_color(*WARM_BEIGE)
                self.set_line_width(0.4)
                self.rect(25, y, self.w - 50, illo_h, "D")
                self.set_font("body", "I", 7)
                self.set_text_color(*SOFT_GREY)
                self.set_xy(25, y + illo_h / 2 - 3)
                self.cell(self.w - 50, 6, f"[ {illustration_hint} ]", align="C")
                y += illo_h + 3

        # Choices
        if choices:
            self.set_text_color(*DARK_BROWN)
            y += 1
            for label, dest in choices:
                self.set_fill_color(*WARM_BEIGE)
                self.rect(20, y, self.w - 40, 6.5, "F")
                self.set_font("body", "B", 7.5)
                self.set_xy(22, y + 1)
                self.cell(self.w - 70, 4.5, strip_emoji(label))
                self.set_font("body", "I", 7)
                self.set_xy(self.w - 50, y + 1)
                self.cell(28, 4.5, f">> Page {dest}", align="R")
                y += 8

    def _cyoa_about_page(self):
        """About this story -- for parents."""
        self.add_page()
        self._bg()
        self._header_bar("About This Story (for parents)")

        y = 26
        sections = [
            ("What is this?",
             "Jem's Big Day Out is a choose-your-own-adventure story about "
             "a party. Your child picks what Jem does at each point. The story "
             "branches based on their choices -- but every path reaches the ending."),
            ("How to read it",
             "Read aloud. At each choice point, let your child pick. Don't steer "
             "them. There are no wrong answers. If they want to go back and try "
             "a different path afterwards, that's great -- it means they're curious."),
            ("What it reveals",
             "Each choice maps to a developmental skill: emotion labelling, "
             "energy reading, social initiation, empathy, perspective-taking. You "
             "don't need to score anything. Just notice which choices feel easy "
             "for your child and which make them pause. That's your data."),
            ("The battery icon",
             "Jem has a battery that fills and drains throughout the story. This "
             "is the Social Gravity spectrum made visible. Children who understand "
             "the battery understand their own energy -- that's the whole point."),
            ("Two entry points",
             "Page 2 asks: does Jem love crowds, or prefer quiet? Pick the one "
             "that matches your child right now. Both paths are equally valid. "
             "The story stretches your child gently in the opposite direction."),
        ]
        for heading, body in sections:
            self.set_font("heading", "B", 8)
            self.set_text_color(*DARK_BROWN)
            self.set_xy(14, y)
            self.cell(self.w - 28, 4, heading)
            y += 4.5
            self.set_font("body", "", 7)
            self.set_xy(14, y)
            self.multi_cell(self.w - 28, 3.3, body, align="L")
            y = self.get_y() + 2.5

    def cyoa_section(self):
        """Build the full Jem's Big Day Out CYOA section."""
        self._story_divider()
        self._cyoa_about_page()

        # PAGE 1 -- Title
        self.add_page()
        self._bg(TAN)
        self.set_font("heading", "B", 24)
        self.set_text_color(*DARK_BROWN)
        self.set_xy(0, 30)
        self.cell(self.w, 10, "Jem's Big Day Out", align="C")
        self.set_font("body", "I", 11)
        self.set_xy(0, 48)
        self.cell(self.w, 6, "A story where YOU choose what happens", align="C")
        # Illustration placeholder
        self.set_draw_color(*WARM_BEIGE)
        self.set_line_width(0.5)
        self.rect(40, 62, self.w - 80, 55, "D")
        self.set_font("body", "I", 9)
        self.set_text_color(*SOFT_GREY)
        self.set_xy(40, 85)
        self.cell(self.w - 80, 8, "[ Jem at a garden gate, balloons visible ]", align="C")
        self.set_font("body", "", 7)
        self.set_text_color(*DARK_BROWN)
        self.set_xy(0, self.h - 12)
        self.cell(self.w, 4, "Page 1", align="C")

        # PAGE 2 -- Entry Point
        self._cyoa_page(2, "The Invitation", [
            ("", "Jem has been invited to a party!"),
            ("", "But here's the thing about Jem..."),
        ], [
            ("Does Jem LOVE being with lots of people?", 3),
            ("Does Jem like having quiet time on their own?", 5),
        ], "Two pictures of Jem -- one with friends, one playing alone. Both happy.")

        # PAGE 3 -- Cohesive Jem arrives
        self._cyoa_page(3, "Cohesive Jem Arrives", [
            ("", "Jem LOVES being with people!"),
            ("", "Jem runs through the gate. There are children everywhere! "
             "Music is playing! Someone is laughing!"),
            ("B", "How does Jem feel right now?"),
        ], [
            ('"Jem feels EXCITED!"', 7),
            ('"Jem feels SO HAPPY!"', 7),
        ], "Jem running in, arms wide, huge grin. Battery icon showing FULL.",
           level_label="Level 1")

        # PAGE 4 -- Cohesive Jem plays
        self._cyoa_page(4, "Jem Plays With Everyone", [
            ("", "Jem plays with EVERYONE."),
            ("", "Jem dances! Jem chases! Jem shares snacks!"),
            ("", "But look at Jem's battery... [HALF]"),
            ("B", "Look at Jem's face. Can you see anything different?"),
        ], [
            ('"Jem still looks happy!"', 8),
            ('"Jem looks a bit tired."', 9),
        ], "Jem still playing but shoulders dropped. Battery HALF.",
           level_label="Level 1")

        # PAGE 5 -- Independent Jem arrives
        self._cyoa_page(5, "Independent Jem Arrives", [
            ("", "Jem likes having quiet time."),
            ("", "Jem walks through the gate. It's... LOUD. "
             "There are children everywhere. Music is playing."),
            ("B", "How does Jem feel right now?"),
        ], [
            ('"Jem feels a bit worried."', 6),
            ('"Jem feels not sure."', 6),
        ], "Jem standing just inside the gate, looking at the party. Battery FULL with question mark.",
           level_label="Level 1")

        # PAGE 6 -- Independent Jem watches
        self._cyoa_page(6, "Jem Watches", [
            ("", "Jem stands near the fence and watches."),
            ("", "Over there -- two children are building something with blocks. "
             "They look like they're having fun. But there's no space to sit down."),
            ("", "Over HERE -- one child is drawing, all on their own. "
             "There's space next to them."),
            ("B", "Which one does Jem look at?"),
        ], [
            ("The children with blocks.", 9),
            ("The child who's drawing.", 8),
        ], "Jem watching from the side. Two scenes: lively group, single child drawing.",
           level_label="Level 1")

        # PAGE 7 -- The party is great!
        self._cyoa_page(7, "The Party Is Great!", [
            ("", "The party is SO much fun!"),
            ("", "Jem has been playing for a long, long time. "
             "Jem's legs are tired. Jem's voice is quiet."),
            ("", "Jem's battery is nearly empty. [LOW]"),
            ("B", "What should Jem do?"),
        ], [
            ('"Keep playing! The party is fun!"', 10),
            ('"Find somewhere quiet to sit."', 11),
            ('"Find their favourite toy and hold it."', 11),
        ], "Jem mid-activity but visibly running low -- yawning, holding own arm.",
           level_label="Level 2")

        # PAGE 8 -- Parallel play
        self._cyoa_page(8, "Side by Side", [
            ("", "Jem sits next to the child who's drawing."),
            ("", "They don't talk. They just draw. Side by side."),
            ("", "It feels... nice."),
            ("B", "How does Jem feel now?"),
        ], [
            ('"Calm."', 11),
            ('"Jem wants to see what they\'re drawing."', 12),
            ('"Jem wants to go play with the big group."', 7),
        ], "Jem and another child sitting side by side, each drawing. Small smiles.",
           level_label="Level 2")

        # PAGE 9 -- Joining in
        self._cyoa_page(9, "Joining In", [
            ("", "Jem walks toward the children with blocks."),
            ("", "They're building a tall tower. It looks wobbly!"),
            ("", "Jem wants to help. But everyone is busy."),
            ("B", "What does Jem do?"),
        ], [
            ('"Jem says: Can I play?"', 12),
            ('"Jem watches and waits."', 8),
            ('"Jem just picks up a block and starts building."', 13),
        ], "Three children building, backs partially turned. Jem standing nearby.",
           level_label="Level 2")

        # PAGE 10 -- Kept playing (consequence)
        self._cyoa_page(10, "Battery Empty", [
            ("", "Jem kept playing... and playing... and playing."),
            ("", "Then Jem's battery ran completely empty. [EMPTY -- RED]"),
            ("", "Jem sat down on the ground and started crying. "
             "The party didn't feel fun anymore."),
            ("B", "What happened?"),
        ], [
            ('"Jem played too long."', 11),
            ('"Jem\'s battery ran out."', 11),
            ("Want to try again? Go back!", 7),
        ], "Jem on the ground, tears. Other children looking confused/concerned.",
           level_label="Level 3")

        # PAGE 11 -- Quiet moment
        self._cyoa_page(11, "The Quiet Spot", [
            ("", "Jem found a quiet spot."),
            ("", "Jem sat down. Jem took a big breath."),
            ("", "Slowly... slowly... the battery started filling up again. [LOW to HALF]"),
            ("", 'A grown-up came over. "Are you OK, Jem?"'),
            ("B", "What does Jem say?"),
        ], [
            ('"I\'m OK. I just needed a rest."', 14),
            ('"I don\'t know."', 14),
            ('"I\'m ready to play again!"', 15),
        ], "Jem sitting against a wall, breathing, eyes closed. Battery rising. Kind adult nearby.",
           level_label="Level 3")

        # PAGE 12 -- Making a connection
        self._cyoa_page(12, "Making a Connection", [
            ("", "Jem is with another child now."),
            ("", 'They\'re drawing together. The other child draws a cat.'),
            ("", '"Look at my cat!" they say.'),
            ("B", "What does Jem do?"),
        ], [
            ('"Jem looks at the cat and says: I like it!"', 14),
            ('"Jem draws a cat too -- right next to it."', 14),
            ('"Jem keeps drawing their own picture."', 13),
        ], "Two children at a table. One holding up a drawing proudly.",
           level_label="Level 3")

        # PAGE 13 -- On your own at the party
        self._cyoa_page(13, "On Your Own", [
            ("", "Jem is on their own at the party."),
            ("", "Everyone else is playing together. Jem is playing alone."),
            ("", "Jem is building something. It's getting really good."),
            ("B", "How does Jem feel?"),
        ], [
            ('"Happy! Jem likes building on their own."', 15),
            ('"A bit lonely. Jem wants someone to see."', 14),
            ('"Jem doesn\'t know yet."', 14),
        ], "Jem alone at a table, building something impressive. Others in background.",
           level_label="Level 3")

        # PAGE 14 -- The middle of the party
        self._cyoa_page(14, "Sam and Alex", [
            ("", "Jem's battery is about half full now. [HALF]"),
            ("", "Two friends come over. Sam and Alex."),
            ("", 'Sam says: "Come and play chase!" '
             'Alex says: "Come and draw with me."'),
            ("", "They both want Jem. But Jem can't do both."),
            ("B", "What does Jem do?"),
        ], [
            ('"Go with Sam -- chase sounds fun!"', 16),
            ('"Go with Alex -- drawing is quieter."', 16),
            ('"Ask if we can all do something together."', 17),
            ('"Can I choose in a minute?"', 17),
        ], "Jem in the middle. Sam bouncing/energetic. Alex sitting with crayons.",
           level_label="Level 4")

        # PAGE 15 -- Jem has energy again
        self._cyoa_page(15, "Where Does Jem Go?", [
            ("", "Jem's battery is filling up! [Nearly FULL]"),
            ("", "The party is still going. There are so many things happening."),
            ("", "Over there: a big loud game with LOTS of children."),
            ("", "Over here: two children doing a puzzle."),
            ("", "In the corner: one child sitting alone, looking down."),
            ("B", "Where does Jem go?"),
        ], [
            ("The big loud game!", 16),
            ("The puzzle with two children.", 17),
            ("The child sitting alone.", 18),
        ], "Wide shot of the party. Three scenes at different energy levels. Child alone in corner.",
           level_label="Level 4")

        # PAGE 16 -- Big energy
        self._cyoa_page(16, "The Big Game", [
            ("", "Jem joins the big game!"),
            ("", "Everyone is running and laughing."),
            ("", "But one child trips and falls. They start crying."),
            ("", "Nobody else stops. They keep running."),
            ("B", "What does Jem do?"),
        ], [
            ('"Keep running! The game is fun!"', 19),
            ('"Stop and check on the child."', 18),
            ('"Shout: STOP! Someone fell!"', 18),
        ], "Chaotic running game. One child on the ground crying. Others still running.",
           level_label="Level 4")

        # PAGE 17 -- Finding the middle
        self._cyoa_page(17, "They Can't Agree", [
            ("", "Jem is with a small group now."),
            ("", "They can't agree on what to do."),
            ("", "One child wants to play outside. One child wants to stay inside."),
            ("", "They're both getting upset."),
            ("B", "What does Jem say?"),
        ], [
            ('"What about playing near the door? A bit of both!"', 19),
            ('"Why do you want to go outside? Why do you want to stay in?"', 19),
            ('"I don\'t know how to fix this."', 19),
        ], "Two children facing each other, arms crossed. Jem in between them.",
           level_label="Level 4")

        # PAGE 18 -- Someone needs help
        self._cyoa_page(18, "Someone Needs Help", [
            ("", "Jem goes to the child who's sitting alone."),
            ("", '"Are you OK?" says Jem.'),
            ("", 'The child says: "My battery ran out. I want to go home. '
             'But the party isn\'t finished."'),
            ("B", "What does Jem do?"),
        ], [
            ('"Jem sits next to them. They don\'t talk. They just sit."', 20),
            ('"Jem says: That happens to me too."', 20),
            ('"Jem brings them a toy."', 20),
            ('"Jem says: I\'ll come back and check on you later."', 20),
        ], "Jem crouching next to the sitting child. Sitting child's battery EMPTY. Jem's HALF.",
           level_label="Level 5")

        # PAGE 19 -- Heading home
        self._cyoa_page(19, "Heading Home", [
            ("", "The party is ending."),
            ("", "Everyone is going home."),
            ("", "Jem had a big day. Some parts were fun. Some parts were hard."),
            ("B", "What was the BEST thing about Jem's day?"),
        ], [
            ('"Playing with everyone!"', 20),
            ('"Drawing with one friend."', 20),
            ('"The quiet time."', 20),
            ('"Helping someone."', 20),
        ], "Jem walking out of the gate, looking back. Warm end-of-day feeling. Battery HALF.",
           level_label="Level 5")

        # PAGE 20 -- The ending
        self.add_page()
        self._bg(CREAM)
        self.set_font("heading", "B", 14)
        self.set_text_color(*DARK_BROWN)
        self.set_xy(0, 20)
        self.cell(self.w, 8, "Home", align="C")

        y = 34
        lines = [
            ("", "Jem is home now."),
            ("", ""),
            ("", "Jem's battery isn't full. And it isn't empty."),
            ("", "It's right in the middle."),
            ("", ""),
            ("B", "That's a good place to be."),
            ("", ""),
            ("", "Tomorrow, Jem might choose different things."),
            ("", "And that's OK too."),
        ]
        for style, text in lines:
            if text:
                self.set_font("body", style, 9)
                self.set_xy(0, y)
                self.cell(self.w, 5, text, align="C")
            y += 5.5

        # Illustration placeholder
        self.set_draw_color(*WARM_BEIGE)
        self.set_line_width(0.5)
        self.rect(45, y + 2, self.w - 90, 35, "D")
        self.set_font("body", "I", 8)
        self.set_text_color(*SOFT_GREY)
        self.set_xy(45, y + 16)
        self.cell(self.w - 90, 6, "[ Jem at home, relaxed. Battery exactly half. Peace. ]", align="C")

        self.set_font("heading", "B", 16)
        self.set_text_color(*DARK_BROWN)
        self.set_xy(0, self.h - 28)
        self.cell(self.w, 8, "The End", align="C")

        self.set_font("body", "I", 8)
        self.set_xy(0, self.h - 18)
        self.cell(self.w, 5, "Want to read it again? Go back to page 2 and try different choices!", align="C")

        self.set_font("body", "", 7)
        self.set_xy(0, self.h - 10)
        self.cell(self.w, 4, "Page 20", align="C")

    # ================================================================
    # Back Matter
    # ================================================================

    def skills_index_page(self):
        """Quick index of the 12 spectra and what they mean."""
        self.add_page()
        self._bg()
        self._header_bar("The 12 Spectra")

        spectra = [
            ("Social Gravity", "Independent", "Cohesive", "How much people charge or drain the battery"),
            ("Energy Directionality", "Quiet", "Loud", "Where their energy naturally sits"),
            ("Voltage Sensitivity", "Low reactor", "High reactor", "How big feelings feel inside"),
            ("Threat Response", "Withdrawn", "Assertive", "What happens when they hear 'no'"),
            ("Care Response", "Self-focused", "Other-focused", "How naturally they think of others"),
            ("Risk Tolerance", "Cautious", "Impulsive", "Whether they jump in or watch first"),
            ("Integrity Logic", "Ordered", "Flexible", "How they feel when rules change"),
            ("Mirror Neuron Tuning", "Selective", "Absorbent", "How much they absorb others' feelings"),
            ("Orderliness", "Messy-tolerant", "Tidy-driven", "How they relate to order and chaos"),
            ("Responsibility Threshold", "Deflecting", "Absorbing", "Who they think things are 'about'"),
            ("Loss Sensitivity", "Easy release", "Strong attachment", "How they handle things going away"),
            ("Libido", "Patient", "Urgent", "How they handle wanting and waiting"),
        ]

        y = 25
        # Table header
        col1, col2, col3, col4 = 42, 24, 24, self.w - 24 - 42 - 24 - 24
        self.set_fill_color(*WARM_BEIGE)
        self.rect(12, y, self.w - 24, 5, "F")
        self.set_font("body", "B", 6.5)
        self.set_text_color(*DARK_BROWN)
        self.set_xy(12, y + 0.5)
        self.cell(col1, 4, "Spectrum")
        self.cell(col2, 4, "Low end")
        self.cell(col3, 4, "High end")
        self.cell(col4, 4, "What it means")
        y += 6

        self.set_font("body", "", 6)
        for name, low, high, desc in spectra:
            self.set_xy(12, y)
            self.set_font("body", "B", 6)
            self.cell(col1, 3.5, name)
            self.set_font("body", "", 6)
            self.cell(col2, 3.5, low)
            self.cell(col3, 3.5, high)
            self.cell(col4, 3.5, desc)
            y += 4.5
            if y > self.h - 15:
                break

        y += 4
        self.set_font("body", "I", 6.5)
        self.set_text_color(*SOFT_GREY)
        self.set_xy(12, y)
        self.multi_cell(self.w - 24, 3.3,
            "None of these are good or bad. They're settings -- biological starting points. "
            "Every child sits somewhere on each spectrum. The activities help them build range "
            "in both directions.", align="C")

    def notes_page(self):
        """Blank notes page for parents."""
        self.add_page()
        self._bg()
        self._header_bar("Notes")

        y = 28
        self.set_font("body", "I", 7.5)
        self.set_text_color(*SOFT_GREY)
        self.set_xy(14, y)
        self.cell(self.w - 28, 4,
                  "What did you notice? What surprised you? What do you want to try next?")
        y += 10

        self.set_draw_color(*WARM_BEIGE)
        self.set_line_width(0.2)
        while y < self.h - 15:
            self.line(14, y, self.w - 14, y)
            y += 7

    def colophon_page(self):
        """Edition info / colophon."""
        self.add_page()
        self._bg(TAN)

        self.set_font("heading", "B", 12)
        self.set_text_color(*DARK_BROWN)
        self.set_xy(0, 20)
        self.cell(self.w, 7, "ToddlerOS: Hello World", align="C")
        self.set_font("body", "", 8)
        self.set_xy(0, 30)
        self.cell(self.w, 5, "Edition 1 -- Foundation Level (All 12 Spectra at Cycle 1)", align="C")

        y = 45
        lines = [
            "Draft for review -- April 2026",
            "",
            "Activities and text by the ToddlerOS team",
            "Illustrations: [Pending]",
            "",
            "Fonts: BLAPEROOM (headings) / Gill Sans Infant Std (body)",
            "",
            "Based on the CHISG developmental skills taxonomy",
            "and the ETP (Emotional Terrain Profile) framework.",
            "",
            "This book works as a standalone physical product.",
            "An optional companion app provides invisible assessment",
            "and personalised pathway recommendations.",
            "",
            "No child was scored, ranked, or compared during the making of this book.",
        ]
        self.set_font("body", "", 7.5)
        for line in lines:
            if line:
                self.set_xy(0, y)
                self.cell(self.w, 4, line, align="C")
            y += 5


def build_book():
    pdf = ToddlerOSBook()

    pdf.cover_page()
    pdf.the_dance_page()
    pdf.how_to_use_page()

    themes = [
        dict(
            num=1, name="Near and Far", spectrum="Social Gravity",
            observation="Which was easier -- playing alone or playing together?",
            cycle_label="Side by Side", activity="The Beside Game",
            desc="Sit next to each other. You draw a picture. They draw a picture. No talking needed. Just being side by side.",
            duration="3 min",
            notice="Do they peek at your drawing? Do they move closer or further away? Do they start talking first, or stay quiet? None of these are right or wrong -- they're data about where your child sits on the Social Gravity spectrum today.",
            rescue="If they get up and walk away, that's fine -- they're telling you their battery for closeness is low right now. If they lean in and start chattering, proximity charges them. Both are information.",
            try_lang="We're just drawing near each other. That's enough.",
            avoid_lang="Don't you want to see what I'm drawing?",
            escalate="They showed comfort -- try The Door Game (Cycle 2).",
            simplify="Three minutes is too long. Try one minute of parallel play.",
            stickers=['"I played near someone"', '"I shared my space"', '"I was happy on my own"', '"I chose to play together"'],
            while_prompt="Two Friends Drawing",
            while_desc="Colour two friends doing their own thing -- together.\n\nThey're not talking. They're not sharing. They're just sitting side by side, each doing their own thing. And that's enough.",
        ),
        dict(
            num=2, name="Loud and Quiet", spectrum="Energy Directionality",
            observation="Did they prefer the noisy game or the quiet one?",
            cycle_label="Volume Dial", activity="The Volume Knob",
            desc="Use your hand as a volume dial. Turn it up -- both get louder. Turn it down -- both get quieter. Take turns being the DJ.",
            duration="3 min",
            notice="Which direction do they turn the dial first? Do they prefer being loud or quiet? When you turn it up, do they match eagerly or hesitate?",
            rescue="If they freeze when it gets loud, they're showing you their comfortable volume. If they can't do quiet, that's data about where their energy naturally sits.",
            try_lang="You're the DJ -- you pick the volume!",
            avoid_lang="Shh! That's too loud!",
            escalate="They switched smoothly -- try Library vs Playground (Cycle 2).",
            simplify="Just do the quiet end. Whisper together.",
            stickers=['"I found my loud voice"', '"I found my quiet voice"', '"I changed my volume"', '"I was the DJ"'],
            while_prompt="Loud and Quiet Colours",
            while_desc="Colour a page where half is loud (bright, big colours) and half is quiet (soft, gentle colours).\n\nWhich side do you like best?",
        ),
        dict(
            num=3, name="Big Feelings, Small Feelings", spectrum="Voltage Sensitivity",
            observation="When something went wrong, was the reaction big or small?",
            cycle_label="The Feelings Thermometer", activity="How Big Is It?",
            desc='Draw a thermometer on paper. Point to scenarios: "Your Lego broke -- how big is that feeling?" They point to the level. No right answer.',
            duration="3 min",
            notice="Can they distinguish between sizes of feelings? A child who says everything is either 'fine' or 'the worst thing ever' is showing you they don't have a middle range yet. That's not drama -- it's a calibration gap.",
            rescue='If they can\'t point to a level, you do it first. "When MY Lego broke, I felt THIS big." You\'re teaching them that feelings have sizes by modelling it.',
            try_lang="There's no right answer. How big does it feel TO YOU?",
            avoid_lang="It's not that bad. Don't be silly.",
            escalate="They can name sizes -- try The Feeling Timer (Cycle 2).",
            simplify="Use just two levels: big and small. Binary first.",
            stickers=['"I measured my feeling"', '"I had a small feeling and noticed it"', '"I had a big feeling and named it"', '"I helped teddy with their feelings"'],
            while_prompt="The Feelings Thermometer",
            while_desc="Colour the feelings thermometer. At the bottom it's cool and calm. At the top it's hot and BIG.\n\nStick a sticker where YOUR feeling is right now.",
        ),
        dict(
            num=4, name="Yes and No", spectrum="Threat Response",
            observation="When told no, did they push back or withdraw?",
            cycle_label="Silly Questions", activity="The Silly Question Game",
            desc='Ask absurd yes/no questions: "Should we put socks on our ears?" "Should the cat drive the car?" They practise saying YES and NO without stakes.',
            duration="3 min",
            notice="Do they default to 'yes' to please you? Do they always say 'no' for power? Can they switch between them freely? A child who can only do one is showing you which end of the Threat Response spectrum they sit on.",
            rescue="If they won't play, they might be worried about wrong answers. Make YOUR answers ridiculous first.",
            try_lang="There's no wrong answer -- it's just silly!",
            avoid_lang="Come on, what's the REAL answer?",
            escalate="They switch freely -- try The Two Plates (Cycle 2).",
            simplify="Just do the silly yes questions. Add no later.",
            stickers=['"I said a big YES"', '"I said a clear NO"', '"I laughed at a silly question"', '"I changed my mind"'],
            while_prompt="Traffic Light Faces",
            while_desc="Colour the traffic light faces:\n\nGreen face = YES!\nRed face = NO!\nAmber face = MAYBE...\n\nWhich face do you use most?",
        ),
        dict(
            num=5, name="Mine and Yours", spectrum="Care Response",
            observation="When asked to share, was it easy or hard?",
            cycle_label="What's Mine", activity="The Sorting Game",
            desc="Get 6 items. Three are yours, three are theirs. Sort them into piles. Name ownership before asking to share.",
            duration="3 min",
            notice="Can they tell which is whose? Do they want ALL of them? Do they want NONE? A child who can't sort ownership isn't being greedy -- they haven't developed the concept of 'mine' and 'yours' as separate categories yet.",
            rescue='If they grab everything, don\'t moralize. Say: "OK, these ones are YOURS. Can you find which ones?" You\'re building the cognitive map before expecting the behaviour.',
            try_lang="Let's sort them. These are yours. These are mine.",
            avoid_lang="Don't be greedy! You have to share.",
            escalate="Ownership is clear -- try The Timer Share (Cycle 2).",
            simplify="Two items, not six. Yours and mine. Just two piles.",
            stickers=['"I know what\'s mine"', '"I know what\'s yours"', '"I sorted them all"', '"I found something that\'s ours"'],
            while_prompt="My Things and Your Things",
            while_desc="Draw your favourite things on one side of the page.\n\nDraw things you'd lend to a friend on the other side.\n\nAre any of them the same?",
        ),
        dict(
            num=6, name="Try and Wait", spectrum="Risk Tolerance",
            observation="Did they jump straight in or watch first?",
            cycle_label="Tiny Bravery", activity="One Taste",
            desc='New food. One tiny taste. "You don\'t have to like it. Just let your tongue have a look."',
            duration="2 min",
            notice="Do they need time before the taste? Do they refuse outright? Do they try immediately? A child who watches the food for two minutes before touching it is NOT being difficult -- they're risk-assessing. That's a skill.",
            rescue="If they refuse, that's their Risk Tolerance position today. Don't push. \"Your tongue said not today. Maybe next time.\"",
            try_lang="Just let your tongue have a look.",
            avoid_lang="Just try it! It's nice!",
            escalate="They tasted it -- try The Spy Game (Cycle 2).",
            simplify="They don't have to taste. Just smell it. Just touch it.",
            stickers=['"I tried something new"', '"I thought about it first"', '"I said not today -- and that\'s OK"', '"I was brave"'],
            while_prompt="My Bravery Badge",
            while_desc="Colour in your bravery badge.\n\nIt doesn't matter if you tried the new thing or not. Being brave isn't about doing scary things. It's about NOTICING that something feels scary.\n\nStick your badge on when you're ready.",
        ),
        dict(
            num=7, name="Same and Different", spectrum="Integrity Logic",
            observation="When the rules changed, was that OK or upsetting?",
            cycle_label="The Rules Game", activity="Today's Rules",
            desc='Make silly rules together: "Today, shoes go in the fridge." Follow them for 5 minutes. Then: "OK, rules back to normal."',
            duration="3 min",
            notice="Do they love the silliness or does it distress them? A child who gets anxious when rules break is showing you they sit on the Ordered end. A child who never wants rules back is showing you the Flexible end.",
            rescue='If they get distressed, honour it. "OK, real rules are back. You like things to stay the same. That\'s a good thing to know about yourself."',
            try_lang="We're being silly on purpose. The real rules will come back.",
            avoid_lang="It's just a game! Why are you upset?",
            escalate="They enjoyed rule-bending -- try The Exception Card (Cycle 2).",
            simplify="One silly rule, not several. Keep it tiny.",
            stickers=['"I followed a new rule"', '"I made up a rule"', '"I was OK when it changed back"', '"The silly rule made me laugh"'],
            while_prompt="Spot the Difference",
            while_desc="Two pictures that are almost the same -- but not quite.\n\nCan you find what's different?\n\nCircle the things that changed.",
        ),
        dict(
            num=8, name="Your Feelings, My Feelings", spectrum="Mirror Neuron Tuning",
            observation="When another child cried, did they notice?",
            cycle_label="Whose Feeling Is This?", activity="The Face Game",
            desc='Make faces at each other: happy, sad, angry, silly. "That\'s MY face. What face are YOU making?" Two different faces -- two different feelings.',
            duration="3 min",
            notice="Do they copy your face automatically (absorbent end) or make their own face independently (selective end)? Can they hold a different face to yours? This reveals how strongly they absorb vs. process other people's emotions.",
            rescue="If they can only copy yours, that's important data -- they may be on the Absorbent end. If they ignore yours completely, Selective end. Neither is wrong.",
            try_lang="Your face is yours. My face is mine. They can be different!",
            avoid_lang="Don't pull that face! That's rude.",
            escalate="They held their own face -- try The Feeling Shield (Cycle 2).",
            simplify="Just do happy and sad. Two faces, not four.",
            stickers=['"I made my own face"', '"I spotted someone else\'s feeling"', '"Our faces were different -- and that\'s OK"', '"I tried lots of feelings"'],
            while_prompt="Two Faces",
            while_desc="Draw two faces.\n\nOne face is how YOU feel right now.\nThe other face is how someone you know feels today.\n\nAre they the same or different?",
        ),
        dict(
            num=9, name="Tidy and Messy", spectrum="Orderliness",
            observation="When it was time to tidy up, was that a fight?",
            cycle_label="Exploring Both", activity="The Messy Minute",
            desc='Timer for 60 seconds. Make the BIGGEST mess you can. Timer stops. Look at it. "That\'s a LOT of mess." Then: "Now let\'s tidy for 60 seconds."',
            duration="3 min",
            notice="Which part did they prefer -- making the mess or tidying it? Did they refuse to make the mess (Ordered end)? Did they refuse to tidy (Flexible end)? Could they switch between both?",
            rescue='If they can\'t make a mess: "It\'s OK. Mess is allowed right now." If they can\'t tidy: "We\'re not making it perfect. Just... a bit less messy."',
            try_lang="Both parts are the game. Messy is step one. Tidy is step two.",
            avoid_lang="Look at this mess! You need to learn to be tidier.",
            escalate="They switched freely -- try Tidy YOUR Way (Cycle 2).",
            simplify="Five objects. Two piles: floor and table. That's it.",
            stickers=['"I made a big mess!"', '"I tidied it all up"', '"I went from messy to tidy"', '"The mess was OK for a minute"'],
            while_prompt="Messy Room, Tidy Room",
            while_desc="Colour a room that's messy on one side and tidy on the other.\n\nNeither is right. Neither is wrong. Sometimes mess is where the fun happens. Sometimes tidy is where the calm lives.",
        ),
        dict(
            num=10, name="My Fault, Your Fault", spectrum="Responsibility Threshold",
            observation="When something went wrong, did they blame others or take all the blame?",
            cycle_label="Sorting the Pieces", activity="Whose Bit Was That?",
            desc='After any small incident, sit together. Draw a pie circle. "This bit was yours. This bit was the table being wobbly." Name each piece without blame.',
            duration="5 min",
            notice='Do they take ALL the blame (Absorbing)? Do they accept NONE (Deflecting)? Can they see the picture as having multiple pieces? A child who says "it\'s all my fault" is not humble -- they\'re absorbing responsibility that isn\'t theirs.',
            rescue='If they shut down: "We\'re not looking for villains. We\'re being detectives. What\'s the evidence?"',
            try_lang="Let's find all the pieces. Yours, mine, the wobbly table's.",
            avoid_lang="Whose fault was that?",
            escalate="They can divide the pie -- try Fix Your Bit (Cycle 2).",
            simplify="Just two pieces: yours and not-yours. Skip the detail.",
            stickers=['"I found my bit"', '"I sorted it fairly"', '"I said what I did"', '"Everyone had a piece"'],
            while_prompt="The Pie Chart",
            while_desc="Draw a circle. This is the pie chart of what happened.\n\nColour each piece a different colour.\nWrite who that piece belongs to.\n\nNot everything is one person's piece.",
        ),
        dict(
            num=11, name="Keeping and Letting Go", spectrum="Loss Sensitivity",
            observation="When something was taken away, was it a crisis or no big deal?",
            cycle_label="Things Go Home", activity="Things Go Home",
            desc='Name where each toy "lives." "Teddy lives on the shelf." Now: "Time for things to go home." Each toy goes HOME -- not "away."',
            duration="5 min",
            notice='Does the word "home" change how they feel about putting things away? A child who panics at "tidy up" but is fine with "things go home" is showing you that the loss frame triggers voltage but the security frame doesn\'t.',
            rescue="If they won't let go of something: \"That one isn't going home yet. It can stay a bit longer.\" Don't force the transition.",
            try_lang="Teddy is going home to the shelf. They'll be there tomorrow.",
            avoid_lang="Put it away NOW. You don't need that.",
            escalate="Things go home easily -- try The Library Game (Cycle 2).",
            simplify="One toy. One home. That's the whole activity.",
            stickers=['"Everything has a home"', '"I said goodnight to my things"', '"Teddy went home safely"', '"Going home is not going away"'],
            while_prompt="Where Things Live",
            while_desc="Draw where your favourite things live.\n\nTeddy lives on the bed. Books live on the shelf. Shoes live by the door.\n\nEverything has a home. And home is safe.",
        ),
        dict(
            num=12, name="Wanting and Waiting", spectrum="Libido",
            observation="When they wanted something, could they wait or did they need it now?",
            cycle_label="Naming What I Want", activity="I Want...",
            desc='Before snack or activity: "What do you want? Say it out loud." Every stated want gets acknowledged: "You want X. I heard you." Then: "Can you have it now, or wait?"',
            duration="3 min",
            notice="Can they name the want? Some children act on desire without ever naming it (they just grab). Others can name it but can't wait. Naming is the first skill -- waiting comes later. They're separate things.",
            rescue='If they can\'t name it, name it for them: "You want the biscuit. I can see you want it." You\'re modelling the language of desire recognition.',
            try_lang="You want X. I heard you. That's a real want.",
            avoid_lang="You can't always have what you want.",
            escalate="They can name wants -- try The 30-Second Wait (Cycle 2).",
            simplify="Just name one want. Don't add the waiting question yet.",
            stickers=['"I said what I wanted"', '"Someone heard me"', '"I thought about what I want"', '"Wanting things is OK"'],
            while_prompt="My Waiting Hourglass",
            while_desc="Colour in the hourglass, one grain at a time.\n\nWaiting is hard. But every grain that falls means you're one grain closer.\n\nYou don't have to be patient. You just have to keep colouring.",
        ),
    ]

    for t in themes:
        pdf.activity_spread(
            theme_num=t["num"], theme_name=t["name"], spectrum=t["spectrum"],
            observation=t["observation"], cycle_label=t["cycle_label"],
            activity=t["activity"], desc=t["desc"], duration=t["duration"],
            notice=t["notice"], rescue=t["rescue"],
            try_lang=t["try_lang"], avoid_lang=t["avoid_lang"],
            escalate=t["escalate"], simplify=t["simplify"],
            stickers=t["stickers"],
            while_prompt=t["while_prompt"], while_desc=t["while_desc"],
        )

    pdf.reference_card()

    # CYOA: Jem's Big Day Out
    pdf.cyoa_section()

    # Back matter
    pdf.skills_index_page()
    pdf.notes_page()
    pdf.colophon_page()

    os.makedirs(os.path.dirname(OUTPUT_PATH), exist_ok=True)
    pdf.output(OUTPUT_PATH)
    print(f"Generated: {OUTPUT_PATH}")
    print(f"Pages: {pdf.page}")


if __name__ == "__main__":
    build_book()
