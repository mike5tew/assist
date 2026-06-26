import React from 'react';
import { Typography, Box } from '@mui/material';
import { Shield as ShieldIcon } from '@mui/icons-material';
import BlogPostLayout from './BlogPostLayout';

const ThreatResponse: React.FC = () => (
  <BlogPostLayout
    title="Threat Response"
    subtitle="What your body does when it can't carry any more. ETP spectrum 4 of 12."
    seriesLabel="ETP Deep-Dive · Post 4 of 12"
    date="March 2026"
    heroColor="#c62828"
    heroGradient="linear-gradient(135deg, #b71c1c 0%, #c62828 45%, #ef5350 100%)"
    icon={<ShieldIcon sx={{ fontSize: 'inherit' }} />}
    tags={['ETPs', 'EdTech', 'AIinEducation', 'EmotionalIntelligence', 'TeacherLife', 'ChildDevelopment']}
    featuredImage="/blog-threat-response.png"
    prevPost={{ label: 'Voltage Sensitivity', path: '/blog/voltage-sensitivity' }}
    nextPost={{ label: 'Care Response', path: '/blog/care-response' }}
  >
    <Typography variant="body1" component="p">
      Your system just overloaded. What happens next?
    </Typography>

    <Typography variant="body1" component="p">
      Post 1: Social Gravity — who charges your battery. Post 2: Energy Directionality — where your processing happens. Post 3: Voltage Sensitivity — how much you can carry. This week: what your body does when it can't carry any more.
    </Typography>

    <Typography variant="body1" component="p">
      <strong>Threat Response.</strong> The spectrum nobody chooses but everybody runs.
    </Typography>

    <Typography variant="h2" component="h2">The Spectrum</Typography>

    <Typography variant="body1" component="p">
      <strong>Freeze</strong> — your system shuts down. You go still, quiet, compliant. Not calm. Offline. The lights are on but the processor stopped. People call you "shy" or "well-behaved."
    </Typography>
    <Typography variant="body1" component="p">
      <strong>Fight</strong> — your system surges. You push back, argue, disrupt. Not aggressive. Overloaded. The circuit is arcing. People call you "defiant" or "challenging."
    </Typography>
    <Typography variant="body1" component="p">
      Neither is a choice. Both kept humans alive:
    </Typography>
    <Box component="ul" sx={{ pl: 3, mb: 2, '& li': { mb: 1, color: 'text.secondary', lineHeight: 1.7 } }}>
      <li>Freeze kept prey invisible to predators — stillness meant survival</li>
      <li>Fight kept threats at distance — noise meant "not worth it"</li>
    </Box>

    <Typography variant="h2" component="h2">The Misconception</Typography>

    <Typography variant="body1" component="p">
      We treat Freeze as good behaviour and Fight as bad behaviour. In classrooms, this is catastrophic.
    </Typography>
    <Typography variant="body1" component="p">
      The Freeze child gets praised for sitting quietly while their brain is in shutdown. They're not learning. They're surviving. And we're rewarding the shutdown.
    </Typography>
    <Typography variant="body1" component="p">
      The Fight child gets sanctioned for pushing back while their system screams overload. They're not defiant. They're doing the only thing their wiring allows. And we're punishing the signal.
    </Typography>

    <Typography variant="h2" component="h2">Your Comfort Zone</Typography>

    <Typography variant="h3" component="h3">Freeze:</Typography>
    <Typography variant="body1" component="p">
      You shut down in confrontation. You can't find words when challenged. You've agreed to things you didn't want because fighting felt impossible. People think you're fine. You're not.
    </Typography>

    <Typography variant="h3" component="h3">Fight:</Typography>
    <Typography variant="body1" component="p">
      You escalate before you've thought. Adrenaline speaks first. You've damaged relationships because your system fired before your brain caught up. People think you're angry. You're flooded.
    </Typography>

    <Typography variant="h2" component="h2">The Classroom Cost</Typography>

    <Typography variant="body1" component="p">
      Teacher raises voice → Freeze child goes silent → Teacher thinks "good, they're listening" → Child has checked out entirely. Learning: zero.
    </Typography>
    <Typography variant="body1" component="p">
      Teacher raises voice → Fight child pushes back → Teacher thinks "defiance" → Sanctions escalate → Child's system confirms: this place is a threat. Trust: gone.
    </Typography>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        Same stimulus. Two responses. Both misread. Both damaged.
      </Typography>
    </Box>

    <Typography variant="h2" component="h2">What Works When You See The Setting</Typography>

    <Typography variant="body1" component="p">
      <strong>Old frame:</strong> "Why won't you answer me?" / "Why are you being rude?"
    </Typography>
    <Typography variant="body1" component="p">
      <strong>New frame:</strong> "They froze — break it down, stay close." / "They're fighting — deflect, re-approach."
    </Typography>

    <Typography variant="h3" component="h3">For the Freeze child: breakdown and persistence</Typography>
    <Typography variant="body1" component="p">
      The freeze response is an overwhelm response. The whole thing feels too big, so the system shuts down. Breaking it down removes the trigger — the task shrinks below the threat threshold. Then persist gently. Stay with them. Don't move on. The message is: "This is small enough. You can do this."
    </Typography>
    <Typography variant="body1" component="p">
      Persistence is the part most people miss. Teachers tend to interpret freeze-silence as "they're fine" and move on to the next child. Staying with them is the signal that it's safe to come back online.
    </Typography>

    <Typography variant="h3" component="h3">For the Fight child: distraction and humour</Typography>
    <Typography variant="body1" component="p">
      This is de-escalation by pattern-interrupt rather than head-on confrontation. A child in fight mode has their sympathetic nervous system firing — matching their energy with authority just confirms the threat. Humour and distraction are lateral moves. They break the loop, let the voltage drop, and create a re-entry point.
    </Typography>
    <Typography variant="body1" component="p">
      The one caveat: it has to be genuine. Kids in fight mode have very sharp threat detection. If the humour feels like a tactic, it can escalate.
    </Typography>

    <Typography variant="h3" component="h3">For both: enthusiasm and praise on re-engagement</Typography>
    <Typography variant="body1" component="p">
      This is the piece that turns a managed episode into a learned behaviour. You're reinforcing the push-through, not the answer. The child starts linking "I got stuck and kept going" with a positive outcome. Over time, that's what builds the capacity to tolerate higher voltage — which is the actual goal.
    </Typography>
    <Typography variant="body1" component="p">
      Without this reinforcement you're just firefighting each episode individually. With it, you're expanding the comfort zone. The child learns that the stuck point isn't the end — it's the beginning of something worth celebrating.
    </Typography>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        These aren't interventions that require a 1:1 therapist. They're reads and adjustments that happen in 30 seconds once you know what you're looking at.
      </Typography>
    </Box>

    <Typography variant="h2" component="h2">The Pattern</Typography>

    <Typography variant="body1" component="p">
      Four posts in, and the pattern builds:
    </Typography>
    <Box component="ul" sx={{ pl: 3, mb: 2, '& li': { mb: 1, color: 'text.secondary', lineHeight: 1.7 } }}>
      <li><strong>Social Gravity</strong> → who charges you</li>
      <li><strong>Energy Directionality</strong> → how you process</li>
      <li><strong>Voltage Sensitivity</strong> → how much you can carry</li>
      <li><strong>Threat Response</strong> → what happens when you can't</li>
    </Box>

    <Typography variant="body1" component="p">
      These interact. A Low Threshold, Inward, Freeze child in a loud classroom isn't quiet because they're focused. They tripped three spectra ago and nobody noticed.
    </Typography>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        Once you see the combination, behaviour stops being a discipline problem and starts being an engineering problem.
      </Typography>
    </Box>

    <Typography variant="body1" component="p">
      <strong>Coming next:</strong> Care Response — whose pain do you feel first? Yours, or everyone else's?
    </Typography>
  </BlogPostLayout>
);

export default ThreatResponse;
