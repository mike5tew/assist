import React from 'react';
import { Typography, Box } from '@mui/material';
import { Groups as GroupsIcon } from '@mui/icons-material';
import BlogPostLayout from './BlogPostLayout';

const SocialGravity: React.FC = () => (
  <BlogPostLayout
    title="Social Gravity"
    subtitle="Who charges your battery — and who drains it. ETP spectrum 1 of 12."
    seriesLabel="ETP Deep-Dive · Post 1 of 12"
    date="March 2026"
    heroColor="#1565c0"
    heroGradient="linear-gradient(135deg, #0d47a1 0%, #1565c0 45%, #42a5f5 100%)"
    icon={<GroupsIcon sx={{ fontSize: 'inherit' }} />}
    tags={['ETPs', 'EdTech', 'AIinEducation', 'EmotionalIntelligence', 'LearningDesign', 'ChildDevelopment', 'TeacherLife', 'NeuroDiversity']}
    featuredImage="/blog-social-gravity.jpeg"
    nextPost={{ label: 'Energy Directionality', path: '/blog/energy-directionality' }}
  >
    <Typography variant="body1" component="p">
      Your brain has a battery. Social interaction either charges it or drains it.
    </Typography>

    <Typography variant="body1" component="p">
      This isn't introversion vs. extroversion — that's a blunt instrument. This is Social Gravity: one of 12 biological spectra we call Emotional Trigger Points (ETPs).
    </Typography>

    <Typography variant="h2" component="h2">The Spectrum</Typography>

    <Typography variant="body1" component="p">
      Social Gravity has two poles:
    </Typography>

    <Typography variant="body1" component="p">
      <strong>Independent</strong> — you recharge alone. Company costs voltage.
    </Typography>
    <Typography variant="body1" component="p">
      <strong>Cohesive</strong> — you recharge with others. Solitude costs voltage.
    </Typography>

    <Typography variant="body1" component="p">
      Neither is better. Both kept humans alive:
    </Typography>

    <Box component="ul" sx={{ pl: 3, mb: 2, '& li': { mb: 1, color: 'text.secondary', lineHeight: 1.7 } }}>
      <li>Independent scouts spotted threats the tribe would miss</li>
      <li>Cohesive members held the group together through connection</li>
    </Box>

    <Typography variant="body1" component="p">
      Your default pole is where you feel most yourself.
    </Typography>

    <Typography variant="h2" component="h2">Your Comfort Zone</Typography>

    <Typography variant="h3" component="h3">If you're Independent:</Typography>
    <Box component="ul" sx={{ pl: 3, mb: 2, '& li': { mb: 1, color: 'text.secondary', lineHeight: 1.7 } }}>
      <li>You think best in silence</li>
      <li>You need alone time like oxygen</li>
      <li>You've been called "aloof" or "distant"</li>
    </Box>

    <Typography variant="h3" component="h3">If you're Cohesive:</Typography>
    <Box component="ul" sx={{ pl: 3, mb: 2, '& li': { mb: 1, color: 'text.secondary', lineHeight: 1.7 } }}>
      <li>You think by talking</li>
      <li>You need connection to feel real</li>
      <li>You've been called "needy" or "too much"</li>
    </Box>

    <Typography variant="h2" component="h2">The Cost</Typography>

    <Typography variant="body1" component="p">
      Life happens at both poles.
    </Typography>
    <Typography variant="body1" component="p">
      The Independent parent still needs to show up for family dinner.
      The Cohesive employee still needs to work alone on that report.
    </Typography>
    <Typography variant="body1" component="p">
      If you can only function at your default, you don't have a preference — you have a prison.
    </Typography>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        Range is the goal. Not changing where you sit on the spectrum — expanding how far you can reach.
      </Typography>
    </Box>

    <Typography variant="h2" component="h2">Why This Matters for Teachers</Typography>

    <Typography variant="body1" component="p">
      You have 30 children. 30 different Social Gravity settings. Running simultaneously.
    </Typography>
    <Typography variant="body1" component="p">
      The Independent child who won't join group work isn't defiant — they're depleted. Their social battery is flashing red.
    </Typography>
    <Typography variant="body1" component="p">
      The Cohesive child who can't stop talking isn't disruptive — they're processing out loud. Silence is where they lose the thread.
    </Typography>
    <Typography variant="body1" component="p">
      <strong>Old frame:</strong> "How do I make Jake cooperate?"
    </Typography>
    <Typography variant="body1" component="p">
      <strong>New frame:</strong> "Jake's reading Independent. He needs parallel play with a gradual invitation — not a demand to join."
    </Typography>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        When you see the setting, behaviour stops being moral and starts being data.
      </Typography>
    </Box>

    <Typography variant="h2" component="h2">Why This Matters for AI</Typography>

    <Typography variant="body1" component="p">
      Most educational AI defaults to interactive, collaborative, screen-forward learning — Cohesive territory.
    </Typography>
    <Typography variant="body1" component="p">
      An AI that doesn't know your Social Gravity will:
    </Typography>
    <Box component="ul" sx={{ pl: 3, mb: 2, '& li': { mb: 1, color: 'text.secondary', lineHeight: 1.7 } }}>
      <li>Push group work on the Independent learner and call disengagement "low motivation"</li>
      <li>Leave the Cohesive learner alone with a screen and wonder why they drift</li>
    </Box>

    <Typography variant="body1" component="p">
      An AI that <em>does</em> know it will say:
    </Typography>
    <Box component="ul" sx={{ pl: 3, mb: 2, '& li': { mb: 1, color: 'text.secondary', lineHeight: 1.7 } }}>
      <li>To the Independent child: "Take 5 minutes to think alone. Then we'll talk."</li>
      <li>To the Cohesive child: "Let's work through this together, out loud."</li>
    </Box>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        That's not personalisation by content. It's personalisation by wiring.
      </Typography>
    </Box>

    <Typography variant="h2" component="h2">The Invitation</Typography>

    <Typography variant="body1" component="p">
      This is post 1 of 12. Each ETP is a spectrum. Each has a comfort zone. Each shapes behaviour in ways we've been moralising instead of measuring.
    </Typography>
    <Typography variant="body1" component="p">
      <strong>Coming next:</strong> Energy Directionality — where your processing power goes (inward vs. outward).
      Then Voltage Sensitivity — how much emotional current your system can carry.
    </Typography>
    <Typography variant="body1" component="p">
      ...and eventually the three new ones: Responsibility Threshold, Loss Sensitivity, and Drive (libido in the engineering sense — the energy behind wanting).
    </Typography>
  </BlogPostLayout>
);

export default SocialGravity;
