import React from 'react';
import { Typography, Box } from '@mui/material';
import { Bolt as BoltIcon } from '@mui/icons-material';
import BlogPostLayout from './BlogPostLayout';

const VoltageSensitivity: React.FC = () => (
  <BlogPostLayout
    title="Voltage Sensitivity"
    subtitle="How much emotional current your system can carry before it trips the breaker. ETP spectrum 3 of 12."
    seriesLabel="ETP Deep-Dive · Post 3 of 12"
    date="March 2026"
    heroColor="#f57c00"
    heroGradient="linear-gradient(135deg, #e65100 0%, #f57c00 45%, #ffb74d 100%)"
    icon={<BoltIcon sx={{ fontSize: 'inherit' }} />}
    tags={['ETPs', 'EdTech', 'AIinEducation', 'EmotionalIntelligence', 'LearningDesign', 'ChildDevelopment', 'TeacherLife', 'NeuroDiversity']}
    featuredImage="/blog-voltage-sensitivity.jpg"
    prevPost={{ label: 'Energy Directionality', path: '/blog/energy-directionality' }}
    nextPost={{ label: 'Threat Response', path: '/blog/threat-response' }}
  >
    <Typography variant="body1" component="p">
      Your circuit breaker has a setting. You didn't choose it.
    </Typography>

    <Typography variant="body1" component="p">
      Post 1: Social Gravity — who charges your battery. Post 2: Energy Directionality — where your processing happens. This week: the one that changes everything else.
    </Typography>

    <Typography variant="body1" component="p">
      <strong>Voltage Sensitivity.</strong> How much emotional current your system can carry before it trips the breaker.
    </Typography>

    <Typography variant="h2" component="h2">The Spectrum</Typography>

    <Typography variant="body1" component="p">
      <strong>High Threshold</strong> — you can absorb enormous emotional load before you shut down. Conflict, chaos, pressure — your system stays online longer. People call you "thick-skinned" or "cold."
    </Typography>
    <Typography variant="body1" component="p">
      <strong>Low Threshold</strong> — your breaker trips early. You feel everything at high resolution. Subtlety, atmosphere, tension — you detect it all. People call you "sensitive" or "dramatic."
    </Typography>
    <Typography variant="body1" component="p">
      Neither is a flaw. Both kept humans alive:
    </Typography>
    <Box component="ul" sx={{ pl: 3, mb: 2, '& li': { mb: 1, color: 'text.secondary', lineHeight: 1.7 } }}>
      <li>High Threshold scouts walked into danger without freezing — someone had to</li>
      <li>Low Threshold watchers detected the threat before anyone else could see it — someone had to do that too</li>
    </Box>

    <Typography variant="h2" component="h2">The Misconception</Typography>

    <Typography variant="body1" component="p">
      We treat this as toughness vs weakness. It isn't. It's <em>resolution</em>.
    </Typography>
    <Typography variant="body1" component="p">
      The Low Threshold person isn't weaker — they're running a higher-resolution emotional sensor. They pick up signals others miss. The cost is that noise hits harder too.
    </Typography>
    <Typography variant="body1" component="p">
      The High Threshold person isn't stronger — they're running a noise filter. They can function in chaos. The cost is they miss signals that matter — the friend who needed help, the student who went quiet.
    </Typography>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        High resolution. Low resolution. Both have blind spots.
      </Typography>
    </Box>

    <Typography variant="h2" component="h2">Your Comfort Zone</Typography>

    <Typography variant="h3" component="h3">Low Threshold:</Typography>
    <Typography variant="body1" component="p">
      Crowded rooms drain you. Raised voices feel like voltage spikes. You need recovery time after emotional events that others shake off in minutes. You've been told to "toughen up" your entire life.
    </Typography>

    <Typography variant="h3" component="h3">High Threshold:</Typography>
    <Typography variant="body1" component="p">
      You genuinely don't understand why someone is upset. You've accidentally hurt people because you didn't feel the weight of your words. You've been told you "don't care" when you simply didn't register the signal.
    </Typography>

    <Typography variant="h2" component="h2">The Cost</Typography>

    <Typography variant="body1" component="p">
      The Low Threshold student in a chaotic classroom isn't learning. Their breaker tripped ten minutes ago. They're sitting in emotional darkness and you think they're daydreaming.
    </Typography>
    <Typography variant="body1" component="p">
      The High Threshold student who just hurt someone's feelings isn't cruel. Their sensor didn't pick up the frequency. Punishing them for something they couldn't detect teaches them nothing — except to mask.
    </Typography>

    <Typography variant="h2" component="h2">The Cost of Getting It Wrong</Typography>

    <Typography variant="body1" component="p">
      When we moralise Voltage Sensitivity:
    </Typography>
    <Box component="ul" sx={{ pl: 3, mb: 2, '& li': { mb: 1, color: 'text.secondary', lineHeight: 1.7 } }}>
      <li>Low Threshold child labelled "anxious" or "attention-seeking" → they learn their sensitivity is a defect → they suppress it → they lose their superpower (early detection)</li>
      <li>High Threshold child labelled "cold" or "uncaring" → they learn empathy is performance → they fake it → they never develop genuine emotional awareness</li>
    </Box>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        Both outcomes are avoidable when you see the setting instead of judging the behaviour.
      </Typography>
    </Box>

    <Typography variant="h2" component="h2">Why This Matters for Teachers</Typography>

    <Typography variant="body1" component="p">
      You have 30 children. 30 different Voltage Sensitivity settings. In one room. With one volume level.
    </Typography>
    <Typography variant="body1" component="p">
      The Low Threshold child who flinches when you raise your voice isn't being defiant. Their system just took a voltage spike it wasn't built for.
    </Typography>
    <Typography variant="body1" component="p">
      The High Threshold child who seems unfazed by a stern conversation isn't disrespectful. The signal genuinely didn't register at their threshold.
    </Typography>
    <Typography variant="body1" component="p">
      <strong>Old frame:</strong> "Why are you so upset?" / "Why don't you care?"
    </Typography>
    <Typography variant="body1" component="p">
      <strong>New frame:</strong> "Their threshold is low — I need to lower the voltage." / "Their threshold is high — I need a different frequency, not more volume."
    </Typography>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        Same situation. Completely different response. Because you read the setting.
      </Typography>
    </Box>

    <Typography variant="h2" component="h2">Why This Matters for AI</Typography>

    <Typography variant="body1" component="p">
      Most educational AI has one emotional setting: neutral. It doesn't detect when a student's voltage is rising, and it doesn't adjust when the breaker trips.
    </Typography>
    <Typography variant="body1" component="p">
      An AI that knows your Voltage Sensitivity:
    </Typography>
    <Box component="ul" sx={{ pl: 3, mb: 2, '& li': { mb: 1, color: 'text.secondary', lineHeight: 1.7 } }}>
      <li><strong>For the Low Threshold learner</strong> — detects frustration early, reduces challenge intensity, offers a pause before the breaker trips. Prevents shutdown instead of reacting to it.</li>
      <li><strong>For the High Threshold learner</strong> — doesn't mistake calm for comprehension. Probes deeper. Uses direct, clear feedback instead of subtle cues they'll never detect.</li>
    </Box>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        Not personalisation by difficulty. Personalisation by emotional resolution.
      </Typography>
    </Box>

    <Typography variant="h2" component="h2">The Pattern</Typography>

    <Typography variant="body1" component="p">
      Three posts in, and the pattern is clear:
    </Typography>
    <Box component="ul" sx={{ pl: 3, mb: 2, '& li': { mb: 1, color: 'text.secondary', lineHeight: 1.7 } }}>
      <li><strong>Social Gravity</strong> → who charges you</li>
      <li><strong>Energy Directionality</strong> → how you process</li>
      <li><strong>Voltage Sensitivity</strong> → how much you can carry</li>
    </Box>

    <Typography variant="body1" component="p">
      These aren't labels. They're settings. They interact. A Low Threshold, Independent, Inward child in a noisy classroom isn't misbehaving — they're overloaded on every spectrum simultaneously.
    </Typography>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        Once you see the combination, behaviour stops being a discipline problem and starts being an engineering problem.
      </Typography>
    </Box>

    <Typography variant="body1" component="p">
      <strong>Coming next:</strong> Emotional Recovery — how quickly your system resets after the breaker trips. Because the threshold is only half the story.
    </Typography>
  </BlogPostLayout>
);

export default VoltageSensitivity;
