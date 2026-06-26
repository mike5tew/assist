import React from 'react';
import { Typography, Box } from '@mui/material';
import { Casino as CasinoIcon } from '@mui/icons-material';
import BlogPostLayout from './BlogPostLayout';

const RiskTolerance: React.FC = () => (
  <BlogPostLayout
    title="Risk Tolerance"
    subtitle="There is no angel. There is no demon. There is just configuration. ETP spectrum 6 of 12."
    seriesLabel="ETP Deep-Dive · Post 6 of 12"
    date="May 2026"
    heroColor="#2e7d32"
    heroGradient="linear-gradient(135deg, #1b5e20 0%, #2e7d32 45%, #43a047 100%)"
    icon={<CasinoIcon sx={{ fontSize: 'inherit' }} />}
    tags={['ETPs', 'EdTech', 'AIinEducation', 'EmotionalIntelligence', 'LearningDesign', 'ChildDevelopment', 'TeacherLife', 'NeuroDiversity', 'Leadership']}
    prevPost={{ label: 'Care Response', path: '/blog/care-response' }}
    nextPost={{ label: 'Integrity Logic', path: '/blog/integrity-logic' }}
  >
    <Typography variant="body1" component="p">
      Part of you wants to leap. Part of you won't let you.
    </Typography>

    <Typography variant="body1" component="p">
      Post 1: Social Gravity — who charges your battery. Post 2: Energy Directionality — where your processing happens. Post 3: Voltage Sensitivity — how much you can carry. Post 4: Threat Response — what you do when it's too much. Post 5: Care Response — whose pain you carry.
    </Typography>

    <Typography variant="body1" component="p">
      This week: the one that feels most like a moral verdict.
    </Typography>

    <Typography variant="body1" component="p">
      <strong>Risk Tolerance.</strong> Your system's relationship with the unknown.
    </Typography>

    <Typography variant="h2" component="h2">There Is No Angel and No Demon</Typography>

    <Typography variant="body1" component="p">
      We've been calling those two parts angel and demon for thousands of years.
    </Typography>

    <Typography variant="body1" component="p">
      We were wrong.
    </Typography>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        There is no angel. There is no demon. There is just a system running multiple subroutines, with different settings on different days. The voice telling you to leap is not your "better nature." It is your Risk Tolerance responding to novelty as reward, filtered through today's Current Load. The voice telling you to hold back is not wisdom. It is a cost calculation being run through three days of accumulated depletion. Two voices. One system. No morality. Just configuration.
      </Typography>
    </Box>

    <Typography variant="body1" component="p">
      The angel/demon model asks: which part of you is right?
    </Typography>
    <Typography variant="body1" component="p">
      The ETP model asks: what is generating each signal, and how much of it is stable versus temporary?
    </Typography>
    <Typography variant="body1" component="p">
      That is a fundamentally different question. And it produces a fundamentally different answer.
    </Typography>

    <Typography variant="h2" component="h2">The Spectrum</Typography>

    <Typography variant="body1" component="p">
      <strong>Risk-Averse</strong> — novelty registers as threat. The unfamiliar is expensive before it's tried. Your system holds back, gathers data, tests the edges. Safe is predictable. New is a cost to be justified.
    </Typography>
    <Typography variant="body1" component="p">
      <strong>Risk-Seeking</strong> — novelty registers as reward. The unfamiliar is oxygen. Your system leans in, tries, fails fast, recalibrates. Familiar is stagnant. New is the whole point.
    </Typography>
    <Typography variant="body1" component="p">
      Neither is courage. Neither is cowardice. Both are wiring.
    </Typography>
    <Typography variant="body1" component="p">
      Both kept the tribe alive: the Risk-Averse remembered which caves had already killed someone and which winters had turned bad. The Risk-Seeking went into the next cave anyway — and found where the fresh water moved. You need both. Your system has a default.
    </Typography>
    <Typography variant="body1" component="p">
      One more thing the spectrum reveals: Risk-Averse is running <strong>high voltage on uncertainty</strong>. Every unknown arrives pre-priced as maximum threat, regardless of evidence. Risk-Seeking is running <strong>low voltage on uncertainty</strong> — the reward signal of novelty drowns out the alarm. Same situation. Different circuitry. The same voltage framework we've been building across all six spectra.
    </Typography>

    <Typography variant="h2" component="h2">The Misconception</Typography>

    <Typography variant="body1" component="p">
      We treat Risk Tolerance as character.
    </Typography>
    <Box component="ul" sx={{ pl: 3, mb: 2, '& li': { mb: 1, color: 'text.secondary', lineHeight: 1.7 } }}>
      <li>Risk-seeking = reckless, impulsive, won't listen.</li>
      <li>Risk-averse = timid, unmotivated, afraid of failure.</li>
    </Box>
    <Typography variant="body1" component="p">
      Both diagnoses are wrong in the same way. They confuse a setting with a verdict.
    </Typography>
    <Typography variant="body1" component="p">
      The risk-seeking child who sprints at every new idea isn't undisciplined. Their reward circuit fires on novelty. Caution doesn't suppress the signal — it just delays it until their system finds stimulation elsewhere. You can't punish risk-seeking out of someone's wiring any more than you can punish voltage sensitivity.
    </Typography>
    <Typography variant="body1" component="p">
      The risk-averse child who won't start the new task isn't lazy. Their threat circuit fires on uncertainty. The blank page isn't blank — it's an unquantified hazard. There is a real cost to every "just try it." They've been paying it their whole life, and nobody acknowledged it.
    </Typography>

    <Typography variant="h2" component="h2">The Consequence Anticipation Link</Typography>

    <Typography variant="body1" component="p">
      Last post I introduced consequence anticipation — the skill beneath Care Response: feel the impact, then simulate the outcome.
    </Typography>
    <Typography variant="body1" component="p">
      Risk Tolerance is not the same as consequence blindness.
    </Typography>
    <Typography variant="body1" component="p">
      The risk-seeking child can foresee consequences. They're discounting them at a different rate. The perceived cost of the bad landing is weighed against the reward of the leap — and their system prices the leap higher. That's not ignorance. It's a different ratio. Teaching risk management isn't about installing fear. It's about calibrating the discount rate: helping them see the full weight of a consequence without destroying the thing that makes them leap.
    </Typography>
    <Typography variant="body1" component="p">
      For the risk-averse person, the ratio is inverted. They're over-counting consequences — every unknown has been given maximum threat probability, regardless of evidence. The skill is accurate probability weighting: "How likely is the bad outcome, really? What does it cost if it goes wrong? What does it cost to never try?"
    </Typography>
    <Typography variant="body1" component="p">
      Both need consequence anticipation. At opposite calibrations.
    </Typography>

    <Typography variant="h2" component="h2">Your Comfort Zone</Typography>

    <Typography variant="h3" component="h3">Risk-Averse:</Typography>
    <Typography variant="body1" component="p">
      You think before you act. You need to see the path before you step onto it. You've been called "overcautious" or "a worrier." You've watched others leap while you were still calculating — sometimes you were right, and sometimes you stood still too long and the moment closed.
    </Typography>

    <Typography variant="h3" component="h3">Risk-Seeking:</Typography>
    <Typography variant="body1" component="p">
      You act before you think. You'd rather ask forgiveness than permission. You've been called "reckless." You've made leaps others wouldn't — sometimes you landed somewhere extraordinary, and sometimes you hit the ground. You don't regret most of them.
    </Typography>

    <Typography variant="h2" component="h2">The Cost</Typography>

    <Typography variant="body1" component="p">
      The risk-averse student in an open-ended project is in crisis. Uncertainty is the assignment. There is no right answer, so every attempt feels like the wrong one. "Just start" is the most unhelpful thing you can say — start what? The blank page is not blank; it's a threat with no defined edge.
    </Typography>
    <Typography variant="body1" component="p">
      The risk-seeking student on a structured, repetitive task is slowly being extinguished. Novelty has gone. The reward circuit isn't firing. What looks like apathy is actually starvation. They're not failing to engage — there's nothing here that their system can engage with.
    </Typography>

    <Typography variant="h2" component="h2">Why This Matters for Teachers</Typography>

    <Typography variant="h3" component="h3">For the risk-averse student:</Typography>
    <Typography variant="body1" component="p">
      Compress the unknowns. This is the <strong>one-inch principle</strong> — make the first step so small and so defined that the threat response doesn't trigger. Not "write an essay" — "write the first sentence." Not "figure out the solution" — "write down what you already know." Reduce the visible unknown. The goal is not to remove caution — it's to give the system a step small enough to step over without tripping the alarm.
    </Typography>

    <Typography variant="h3" component="h3">For the risk-seeking student:</Typography>
    <Typography variant="body1" component="p">
      Inject novelty into the familiar. Gamify. Add time pressure. Change the frame. The content can stay identical — the wrapper moves. "Can you solve this in under two minutes?" does more for a risk-seeking child than any amount of reasoning about why the task matters.
    </Typography>

    <Typography variant="body1" component="p">
      <strong>Old frame:</strong> "Why won't they just try?" / "Why won't they just settle down?"
    </Typography>
    <Typography variant="body1" component="p">
      <strong>New frame:</strong> "Their threat circuit fires on uncertainty — give them a smaller unknown." / "Their reward circuit fires on novelty — give them a moving target."
    </Typography>

    <Typography variant="h2" component="h2">Why This Matters for AI — and humanOS</Typography>

    <Typography variant="body1" component="p">
      Most AI coaching products make the angel/demon error. They respond to the "cautious signal" with reassurance, or to the "bold signal" with encouragement. They pick a side. They play the angel — the calm, wise voice whispering guidance against your worse instincts.
    </Typography>
    <Typography variant="body1" component="p">
      humanOS is built around a different principle.
    </Typography>
    <Typography variant="body1" component="p">
      When someone reports internal conflict around a decision, the question is not: which voice should win? The question is: what is generating each signal, and how much of it is stable configuration versus current state?
    </Typography>
    <Typography variant="body1" component="p">
      It sounds like this:
    </Typography>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        "Your caution around this is being amplified by a Current Load reading that's been elevated for three days. Your baseline Risk Tolerance sits at moderate-seeking — on a calmer week, you'd lean toward this without much friction. The resistance you're feeling right now is not primarily your character. It's your current state. That's worth knowing."
      </Typography>
    </Box>

    <Typography variant="body1" component="p">
      That is not the angel saying "be brave." It is a system readout. It separates the stable profile from the temporary condition. It gives the person something actually useful: information about what is running, so they can decide how much weight to give each signal.
    </Typography>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        Not moral guidance. Configuration clarity. Not two voices debating. One system, and a map of its current state.
      </Typography>
    </Box>

    <Typography variant="h2" component="h2">The Pattern</Typography>

    <Typography variant="body1" component="p">
      Six spectra in:
    </Typography>
    <Box component="ul" sx={{ pl: 3, mb: 2, '& li': { mb: 1, color: 'text.secondary', lineHeight: 1.7 } }}>
      <li><strong>Social Gravity</strong> → who charges you</li>
      <li><strong>Energy Directionality</strong> → how you process</li>
      <li><strong>Voltage Sensitivity</strong> → how much you carry</li>
      <li><strong>Threat Response</strong> → what happens at the limit</li>
      <li><strong>Care Response</strong> → whose pain you absorb</li>
      <li><strong>Risk Tolerance</strong> → how you price the unknown</li>
    </Box>

    <Typography variant="body1" component="p">
      These are not six separate labels. They are six settings running simultaneously. The "internal debate" you experience — the thing that makes certain decisions hard — is not a moral dialogue between two agents. It is the computational output of a multi-setting system, interacting in real time.
    </Typography>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        The combination explains the experience. Every time.
      </Typography>
    </Box>

    <Typography variant="body1" component="p">
      <strong>Coming next:</strong> Integrity Logic. What happens when the right choice and the natural choice aren't the same thing?
    </Typography>
  </BlogPostLayout>
);

export default RiskTolerance;
