import React from 'react';
import { Typography, Box } from '@mui/material';
import { SwapHoriz as SwapIcon } from '@mui/icons-material';
import BlogPostLayout from './BlogPostLayout';

const EnergyDirectionality: React.FC = () => (
  <BlogPostLayout
    title="Energy Directionality"
    subtitle="Where your processing happens — inward or outward. ETP spectrum 2 of 12."
    seriesLabel="ETP Deep-Dive · Post 2 of 12"
    date="March 2026"
    heroColor="#6a1b9a"
    heroGradient="linear-gradient(135deg, #4a148c 0%, #6a1b9a 45%, #ab47bc 100%)"
    icon={<SwapIcon sx={{ fontSize: 'inherit' }} />}
    tags={['ETPs', 'EdTech', 'AIinEducation', 'EmotionalIntelligence', 'LearningDesign', 'ChildDevelopment', 'TeacherLife', 'NeuroDiversity']}
    featuredImage="/blog-energy-directionality.jpg"
    prevPost={{ label: 'Social Gravity', path: '/blog/social-gravity' }}
    nextPost={{ label: 'Voltage Sensitivity', path: '/blog/voltage-sensitivity' }}
  >
    <Typography variant="body1" component="p">
      Your brain is a processor. But it runs in one of two directions.
    </Typography>

    <Typography variant="body1" component="p">
      Last week: Social Gravity — whether people charge or drain your battery. This week: Energy Directionality. The spectrum everyone confuses with it.
    </Typography>

    <Typography variant="body1" component="p">
      They're not the same. Social Gravity is about <em>who</em> charges you. Energy Directionality is about <em>where</em> your processing happens.
    </Typography>

    <Typography variant="h2" component="h2">The Spectrum</Typography>

    <Typography variant="body1" component="p">
      <strong>Inward</strong> — you think, then speak. The world is input; your mind is the workspace.
    </Typography>
    <Typography variant="body1" component="p">
      <strong>Outward</strong> — you speak to think. The world is your workspace; other people are the sounding board.
    </Typography>

    <Typography variant="h2" component="h2">The Confusion</Typography>

    <Typography variant="body1" component="p">
      We collapse this into introvert/extrovert. But you can be:
    </Typography>
    <Box component="ul" sx={{ pl: 3, mb: 2, '& li': { mb: 1, color: 'text.secondary', lineHeight: 1.7 } }}>
      <li><strong>Independent + Outward</strong> — you don't need people to recharge, but you think out loud. You're the person muttering at your desk.</li>
      <li><strong>Cohesive + Inward</strong> — you need people around, but you think silently. You're the quiet one at the party who still needs to be there.</li>
    </Box>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        Most personality models can't see this. They've welded two spectra into one label.
      </Typography>
    </Box>

    <Typography variant="h2" component="h2">Your Comfort Zone</Typography>

    <Typography variant="h3" component="h3">Inward:</Typography>
    <Box component="ul" sx={{ pl: 3, mb: 2, '& li': { mb: 1, color: 'text.secondary', lineHeight: 1.7 } }}>
      <li>You need time before answering. Being put on the spot feels like an ambush.</li>
      <li>Your best ideas come 20 minutes after the meeting.</li>
    </Box>

    <Typography variant="h3" component="h3">Outward:</Typography>
    <Box component="ul" sx={{ pl: 3, mb: 2, '& li': { mb: 1, color: 'text.secondary', lineHeight: 1.7 } }}>
      <li>You need to talk to solve. Silence during brainstorming is suffocation.</li>
      <li>You've said things you hadn't thought yet — out loud.</li>
    </Box>

    <Typography variant="h2" component="h2">The Cost</Typography>

    <Typography variant="body1" component="p">
      The Inward processor in a fast meeting doesn't get heard — the conversation moved on before their processing cycle completed.
    </Typography>
    <Typography variant="body1" component="p">
      The Outward processor in an exam sits in silence wondering why their brain won't start. No one to bounce off. The engine needs external ignition.
    </Typography>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        Neither is failing. Different architectures, environment designed for one.
      </Typography>
    </Box>

    <Typography variant="h2" component="h2">Why This Matters for Teachers</Typography>

    <Typography variant="body1" component="p">
      <strong>Discussion task:</strong> Outward processors light up — this IS their thinking. Inward processors go quiet. You assume disengaged.
    </Typography>
    <Typography variant="body1" component="p">
      <strong>Silent writing:</strong> Inward processors come alive. Outward processors stare at the page. You assume unprepared.
    </Typography>
    <Typography variant="body1" component="p">
      Sophie processes Inward. She needs 90 seconds of think-time before you ask her to speak.
    </Typography>
    <Typography variant="body1" component="p">
      Tom processes Outward. Give him a partner for three minutes. Then he'll write it in one draft.
    </Typography>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        Tiny adjustment. Transformational impact.
      </Typography>
    </Box>

    <Typography variant="h2" component="h2">Why This Matters for AI</Typography>

    <Typography variant="body1" component="p">
      Current AI tutoring assumes Inward processing — present information, wait for typed input. Silent exchange between screen and mind. Paradise for Inward. Dead zone for Outward.
    </Typography>
    <Typography variant="body1" component="p">
      An AI that knows your Energy Directionality gives the Inward learner thinking time and accepts delayed, precise responses. It gives the Outward learner a sounding board — encouraging messy, iterative exploration in real-time.
    </Typography>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        Not personalisation by difficulty. Personalisation by processing architecture.
      </Typography>
    </Box>

    <Typography variant="h2" component="h2">The Pattern</Typography>

    <Typography variant="body1" component="p">
      Every ETP spectrum describes something we already see but mislabel.
    </Typography>
    <Box component="ul" sx={{ pl: 3, mb: 2, '& li': { mb: 1, color: 'text.secondary', lineHeight: 1.7 } }}>
      <li>Quiet ≠ disengaged. It might mean Inward.</li>
      <li>Loud ≠ showing off. It might mean Outward.</li>
    </Box>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        When you see the setting, behaviour stops being moral and starts being data.
      </Typography>
    </Box>

    <Typography variant="body1" component="p">
      <strong>Coming next:</strong> Voltage Sensitivity — how much emotional current your system can carry before it trips the breaker.
    </Typography>
  </BlogPostLayout>
);

export default EnergyDirectionality;
