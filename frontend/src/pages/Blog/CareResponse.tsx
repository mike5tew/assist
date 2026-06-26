import React from 'react';
import { Typography, Box } from '@mui/material';
import { Favorite as FavoriteIcon } from '@mui/icons-material';
import BlogPostLayout from './BlogPostLayout';

const CareResponse: React.FC = () => (
  <BlogPostLayout
    title="Care Response"
    subtitle="Whose pain do you feel first? Yours — or everyone else's? ETP spectrum 5 of 12."
    seriesLabel="ETP Deep-Dive · Post 5 of 12"
    date="April 2026"
    heroColor="#ad1457"
    heroGradient="linear-gradient(135deg, #880e4f 0%, #ad1457 45%, #e91e63 100%)"
    icon={<FavoriteIcon sx={{ fontSize: 'inherit' }} />}
    tags={['ETPs', 'EdTech', 'AIinEducation', 'EmotionalIntelligence', 'TeacherLife', 'ChildDevelopment', 'SEN', 'Empathy']}
    prevPost={{ label: 'Threat Response', path: '/blog/threat-response' }}
    nextPost={{ label: 'Risk Tolerance', path: '/blog/risk-tolerance' }}
  >
    <Typography variant="body1" component="p">
      Whose pain do you feel first? Yours — or everyone else's?
    </Typography>

    <Typography variant="body1" component="p">
      Post 1: Social Gravity — who charges your battery. Post 2: Energy Directionality — where your processing happens. Post 3: Voltage Sensitivity — how much you can carry. Post 4: Threat Response — what you do when it's too much.
    </Typography>

    <Typography variant="body1" component="p">
      This week: the one that gets moralised more than any other.
    </Typography>

    <Typography variant="body1" component="p">
      <strong>Care Response.</strong> How your system allocates emotional resources when someone else is hurting.
    </Typography>

    <Typography variant="h2" component="h2">The Spectrum</Typography>

    <Typography variant="body1" component="p">
      <strong>Detached</strong> — other people's pain registers, but at low volume. You process it logically. You can function in a crisis because their distress doesn't flood your system. People call you "cold" or say you "don't care."
    </Typography>
    <Typography variant="body1" component="p">
      <strong>Nurturing</strong> — other people's pain registers at full volume. You feel it in your body before your brain catches up. You can't not respond. People call you "sensitive" or "too soft."
    </Typography>
    <Typography variant="body1" component="p">
      Neither is a defect. Both kept humans alive:
    </Typography>
    <Box component="ul" sx={{ pl: 3, mb: 2, '& li': { mb: 1, color: 'text.secondary', lineHeight: 1.7 } }}>
      <li>Detached members made decisions under pressure — someone had to think while the village was screaming</li>
      <li>Nurturing members detected suffering before anyone asked for help — someone had to notice the child who went quiet</li>
    </Box>

    <Typography variant="h2" component="h2">The Misconception</Typography>

    <Typography variant="body1" component="p">
      We treat this as kindness vs selfishness. It isn't. It's allocation.
    </Typography>
    <Typography variant="body1" component="p">
      The Detached person isn't heartless — they're running low-resolution empathy so they can still function. They care, but feeling the full weight of your pain would crash their system. So they don't. That's not moral failure. That's resource management.
    </Typography>
    <Typography variant="body1" component="p">
      The Nurturing person isn't a saint — they're running high-resolution empathy. They absorb everything. The cost is they carry weight that isn't theirs. They burn out not from effort but from feeling. And we praise them for it right up until they break.
    </Typography>
    <Typography variant="body1" component="p">
      Low resolution. High resolution. Both have blind spots.
    </Typography>

    <Typography variant="h2" component="h2">Your Comfort Zone</Typography>

    <Typography variant="h3" component="h3">Detached:</Typography>
    <Typography variant="body1" component="p">
      You've been told you're "cold" or "unfeeling" when you genuinely didn't register the weight of someone's situation. You handle emergencies well but get blindsided by: "How could you not see I was struggling?"
    </Typography>

    <Typography variant="h3" component="h3">Nurturing:</Typography>
    <Typography variant="body1" component="p">
      You feel other people's pain in your chest. You take on problems that aren't yours. You've been exhausted by a conversation that the other person walked away from feeling fine. You've heard "stop worrying about everyone else" your whole life.
    </Typography>

    <Typography variant="h2" component="h2">The Cost</Typography>

    <Typography variant="body1" component="p">
      This is where it gets dangerous — because we moralise Care Response harder than almost any other spectrum.
    </Typography>
    <Typography variant="body1" component="p">
      The Detached child who doesn't cry when a friend is upset isn't cruel. Their sensor didn't pick up the signal at the frequency it was sent. Telling them "you should feel bad" teaches them that empathy is performance. They learn to fake it. They never develop genuine awareness because we punished the absence instead of teaching the skill.
    </Typography>
    <Typography variant="body1" component="p">
      The Nurturing child who falls apart when someone else is in trouble isn't weak. Their system absorbed a voltage spike meant for someone else. Telling them "it's not your problem" doesn't help — they know it's not their problem. They just can't stop feeling it. They need help sorting which feelings are theirs and which they caught from the room.
    </Typography>

    <Typography variant="h2" component="h2">The 2×2 Nobody Talks About</Typography>

    <Typography variant="body1" component="p">
      Care Response doesn't exist in isolation. Pair it with Threat Response and you get four very different humans:
    </Typography>

    <Box component="ul" sx={{ pl: 3, mb: 2, '& li': { mb: 1, color: 'text.secondary', lineHeight: 1.7 } }}>
      <li><strong>Passive + Detached = Apathetic.</strong> Checked out. Neither their pain nor yours moves them to action.</li>
      <li><strong>Passive + Nurturing = Gentle Caregiver.</strong> Absorbs your pain but won't push back. Burns out quietly.</li>
      <li><strong>Aggressive + Detached = Predatory.</strong> Acts without feeling the impact. The one we're all afraid of.</li>
      <li><strong>Aggressive + Nurturing = Protective Parent.</strong> Feels everything and fights for it. The one who runs toward the screaming.</li>
    </Box>

    <Typography variant="body1" component="p">
      But care alone isn't enough. A child who hits isn't missing empathy — they're missing the ability to simulate: "If I do this, what happens next?" Overprotective parents don't lack care — they lack accurate risk prediction. The skill beneath both is <strong>consequence anticipation</strong>. Feel the impact. Then predict the outcome. That's the sequence.
    </Typography>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        The goal isn't to eliminate aggression. It's to pair it with prediction — so the child can feel the weight and foresee the landing. Aggression without prediction is dangerous. Aggression with prediction is protection.
      </Typography>
    </Box>

    <Typography variant="h2" component="h2">Why This Matters for Teachers</Typography>

    <Typography variant="body1" component="p">
      You have 30 children. 30 different Care Response settings. And a classroom full of emotional spills.
    </Typography>

    <Typography variant="h3" component="h3">For the Detached child:</Typography>
    <Typography variant="body1" component="p">
      They seem unbothered when they hurt someone's feelings. They're not defiant. The signal didn't register. Saying "How would YOU feel?" doesn't work — they literally cannot simulate it at the resolution you're expecting.
    </Typography>
    <Typography variant="body1" component="p">
      Instead: use third-person perspective. Stories, characters, puppets. "How do you think that character felt?" Remove the personal charge. Let them build the skill without the shame.
    </Typography>

    <Typography variant="h3" component="h3">For the Nurturing child:</Typography>
    <Typography variant="body1" component="p">
      They keep trying to fix everyone's problems. They're not interfering. Their system won't let them not respond. The pain is loud and they don't have a volume knob.
    </Typography>
    <Typography variant="body1" component="p">
      Instead: teach sorting. "This feeling is yours. This feeling is theirs. This feeling is from the room. Let's put them in separate boxes." Give them the volume knob they don't have yet.
    </Typography>

    <Typography variant="body1" component="p">
      <strong>Old frame:</strong> "You don't care" / "You care too much."
    </Typography>
    <Typography variant="body1" component="p">
      <strong>New frame:</strong> "Your resolution is low — let's build a signal amplifier." / "Your resolution is high — let's build a filter."
    </Typography>

    <Typography variant="h2" component="h2">Why This Matters for AI</Typography>

    <Typography variant="body1" component="p">
      Most educational AI treats emotional support as one-size-fits-all: encouraging, warm, patient. That works for the Nurturing learner. For the Detached learner, it feels like noise.
    </Typography>
    <Typography variant="body1" component="p">
      An AI that knows your Care Response setting:
    </Typography>
    <Typography variant="body1" component="p">
      <strong>For the Detached learner</strong> — uses direct, factual feedback. Doesn't try to make them feel something. Instead, builds perspective-taking through structured scenarios: "What would happen if this character did X?"
    </Typography>
    <Typography variant="body1" component="p">
      <strong>For the Nurturing learner</strong> — detects when they're absorbing group frustration rather than processing their own work. Intervenes with: "That feeling might not be yours. Let's focus on your task for the next five minutes."
    </Typography>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        Not personalisation by warmth. Personalisation by emotional resolution.
      </Typography>
    </Box>

    <Typography variant="h2" component="h2">The Pattern</Typography>

    <Typography variant="body1" component="p">
      Five spectra in, the picture builds:
    </Typography>
    <Box component="ul" sx={{ pl: 3, mb: 2, '& li': { mb: 1, color: 'text.secondary', lineHeight: 1.7 } }}>
      <li><strong>Social Gravity</strong> → who charges you</li>
      <li><strong>Energy Directionality</strong> → how you process</li>
      <li><strong>Voltage Sensitivity</strong> → how much you carry</li>
      <li><strong>Threat Response</strong> → what you do when it's too much</li>
      <li><strong>Care Response</strong> → whose pain you carry</li>
    </Box>

    <Typography variant="body1" component="p">
      These aren't five separate labels. They're five settings running simultaneously. A Nurturing, Low Threshold, Passive child in a classroom argument isn't "overreacting" — they're absorbing everyone's voltage, at maximum resolution, with no way to discharge it. Of course they're crying. The engineering explains the behaviour.
    </Typography>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        Once you see the combination, behaviour stops being a character judgement and starts being a circuit diagram.
      </Typography>
    </Box>

    <Typography variant="body1" component="p">
      <strong>Coming next:</strong> Risk Tolerance — and the skill that makes every other spectrum safe: consequence anticipation.
    </Typography>
  </BlogPostLayout>
);

export default CareResponse;
