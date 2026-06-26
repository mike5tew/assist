import React from 'react';
import { Typography, Box } from '@mui/material';
import { FlightTakeoff as FlightTakeoffIcon } from '@mui/icons-material';
import BlogPostLayout from './BlogPostLayout';

const PilotStrength: React.FC = () => (
  <BlogPostLayout
    title="Pilot Strength"
    subtitle="What happens when you know exactly what you should do — and can't make yourself do it? ETP spectrum 10 of 12."
    seriesLabel="ETP Deep-Dive · Post 10 of 12"
    date="June 2026"
    heroColor="#5d4037"
    heroGradient="linear-gradient(135deg, #3e2723 0%, #5d4037 45%, #8d6e63 100%)"
    icon={<FlightTakeoffIcon sx={{ fontSize: 'inherit' }} />}
    tags={['ETPs', 'EdTech', 'AIinEducation', 'EmotionalIntelligence', 'LearningDesign', 'ExecutiveFunction', 'TeacherLife', 'NeuroDiversity', 'Leadership']}
    prevPost={{ label: 'Orderliness', path: '/blog/orderliness' }}
    nextPost={{ label: 'Current Load', path: '/blog/current-load' }}
  >
    <Typography variant="body1" component="p">
      You know what you should do. The plan is clear. The intention is real.
    </Typography>

    <Typography variant="body1" component="p">
      And then you do not do it.
    </Typography>

    <Typography variant="body1" component="p">
      That gap — between knowing and doing — is where most of personal development quietly fails. We call it weakness, or laziness, or lack of willpower. It is none of those. It is a spectrum, and it has a name.
    </Typography>

    <Typography variant="body1" component="p">
      <strong>Pilot Strength.</strong> Your system's capacity to override an immediate impulse in favour of a deliberated, longer-term response — the strength of the layer that flies the plane rather than letting it fly itself.
    </Typography>

    <Typography variant="h2" component="h2">The Spectrum</Typography>

    <Typography variant="body1" component="p">
      <strong>High Pilot Strength</strong> — the deliberative layer engages readily. You choose against instinct. You delay, override, redirect. You are called disciplined, focused, self-controlled. Also: overthinking, slow to act, stuck in your own head.
    </Typography>
    <Typography variant="body1" component="p">
      <strong>Low Pilot Strength</strong> — the automatic runs first. Habits fire. Instinct moves before deliberation arrives. You are called impulsive, inconsistent, unreliable. Also: fast, decisive, natural, unhesitating when hesitation would be fatal.
    </Typography>

    <Typography variant="body1" component="p">
      Neither is discipline. Neither is recklessness. Both are wiring.
    </Typography>

    <Typography variant="body1" component="p">
      High-Pilot restraint kept the tribe from acting on every threat and every impulse — the deferred gratification that let people plant rather than only forage. Low-Pilot speed moved when deliberation would have been lethal delay — the fast-twitch response that survives the ambush. You need both. Your system has a ratio.
    </Typography>

    <Typography variant="h2" component="h2">The Brakes and the Engine</Typography>

    <Typography variant="body1" component="p">
      Think of an experienced driver.
    </Typography>

    <Typography variant="body1" component="p">
      Most of the drive is automatic — gear changes, mirror checks, the constant micro-adjustments of brake and steering. The deliberative layer is not running the whole time. It is in reserve.
    </Typography>

    <Typography variant="body1" component="p">
      Pilot Strength is what fires when the automatic fails. The child who runs into the road. The patch of black ice. The vehicle braking hard in front. That is the moment the deliberative override has to arrive — fast, decisive, and present.
    </Typography>

    <Typography variant="body1" component="p">
      And here is the part people miss: if the driver consciously deliberated every action through the entire ordinary drive, the resource would be exhausted before the emergency ever arrived. The automatic layer is not laziness. It is conservation. Your automatic layer is the engine. Pilot Strength is the brakes. You develop the override so that it is there when you actually need it — not so it runs every second.
    </Typography>

    <Typography variant="h2" component="h2">The Internal Voice</Typography>

    <Typography variant="body1" component="p">
      "I know what I should do. I just cannot make myself do it."
    </Typography>

    <Typography variant="body1" component="p">
      That is the deliberative layer failing to engage in time. And under high Current Load, Pilot Strength is the first thing to deplete. What you could normally choose, you cannot today — not because your values changed, but because the resource is running low.
    </Typography>

    <Typography variant="body1" component="p">
      This is why behaviour deteriorates at 3pm, and at the end of term, and at the end of a hard week. Not attitude. Depletion. The executive layer goes offline first, and the autopilot spectra keep running without it.
    </Typography>

    <Typography variant="h2" component="h2">Why This Matters for Teachers</Typography>

    <Typography variant="body1" component="p">
      The student with low Pilot Strength does not need to try harder. They already know what they should do. The deliberative layer is not arriving in time to carry it.
    </Typography>

    <Typography variant="body1" component="p">
      Give them the pause they cannot find internally. Wait-time. A countdown. A physical cue agreed in advance. External structure that stands in for the internal override until the internal one develops. The scaffold is not a crutch — it is borrowed executive capacity during construction.
    </Typography>

    <Typography variant="body1" component="p">
      The student with high Pilot Strength needs the opposite: resistance. Without something worth overriding for, the strength has nothing to work on. They can look uninvolved, disengaged, flat. They are not switched off. They are idling — waiting for something that earns the engagement.
    </Typography>

    <Typography variant="body1" component="p">
      <strong>Old frame:</strong> "They know what to do — they're just not doing it. That's a choice."
    </Typography>
    <Typography variant="body1" component="p">
      <strong>New frame:</strong> "The deliberative layer isn't arriving in time. Is it absent, or is it depleted? Those need different responses."
    </Typography>

    <Typography variant="h2" component="h2">Why This Matters for AI — and humanOS</Typography>

    <Typography variant="body1" component="p">
      An AI teaching assistant built on this model does not read "won't" where the truth is "can't yet." It reads Pilot Strength alongside Current Load, and separates a capacity issue from a character verdict. And humanOS does the same for the adult running it:
    </Typography>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        "The feeling that you know what to do and cannot make yourself do it is not failure. It is the gauge showing low. This is recoverable."
      </Typography>
    </Box>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        Not a motivational prompt. A state reading.
      </Typography>
    </Box>

    <Typography variant="h2" component="h2">The Next Inch</Typography>

    <Typography variant="body1" component="p">
      Today, notice one moment when you knew what to do and did not do it. Do not judge it. Just locate it.
    </Typography>

    <Typography variant="body1" component="p">
      Then ask one question: was the deliberative layer <em>absent</em> — or <em>depleted</em>? If it was absent, the work is building the override through practice and external scaffolding. If it was depleted, the work is restoring buffer before you ask anything more of it. They look identical from the outside. They are different problems, and they need different solutions.
    </Typography>

    <Typography variant="h2" component="h2">The Pattern</Typography>

    <Typography variant="body1" component="p">
      Ten spectra now:
    </Typography>
    <Box component="ul" sx={{ pl: 3, mb: 2, '& li': { mb: 1, color: 'text.secondary', lineHeight: 1.7 } }}>
      <li><strong>Social Gravity</strong> → who charges you</li>
      <li><strong>Energy Directionality</strong> → how you process</li>
      <li><strong>Voltage Sensitivity</strong> → how much you carry</li>
      <li><strong>Threat Response</strong> → what happens at the limit</li>
      <li><strong>Care Response</strong> → whose pain you absorb</li>
      <li><strong>Risk Tolerance</strong> → how you price the unknown</li>
      <li><strong>Integrity Logic</strong> → how you apply rules under pressure</li>
      <li><strong>Mirror Neuron Tuning</strong> → how much of what you feel is from others</li>
      <li><strong>Orderliness</strong> → what the environment needs to be before work can start</li>
      <li><strong>Pilot Strength</strong> → whether the deliberative layer arrives in time</li>
    </Box>

    <Typography variant="body1" component="p">
      Not ten moral categories. Ten settings, running simultaneously, interacting in real time.
    </Typography>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        The automatic layer is the engine. Pilot Strength is the brakes. You don't run the brakes the whole drive — you develop them so they're there when the child runs into the road.
      </Typography>
    </Box>

    <Typography variant="body1" component="p">
      <strong>Coming next:</strong> Current Load. What happens when the buffer runs out — and every other reading changes?
    </Typography>
  </BlogPostLayout>
);

export default PilotStrength;
