import React from 'react';
import { Typography, Box } from '@mui/material';
import { Gavel as GavelIcon } from '@mui/icons-material';
import BlogPostLayout from './BlogPostLayout';

const TheGovernor: React.FC = () => (
  <BlogPostLayout
    title="The Governor"
    subtitle="An AI role that audits both sides of a disagreement — and is built never to declare a winner."
    seriesLabel="Applied ETP · Evidence & Disagreement"
    date="June 2026"
    heroColor="#283593"
    heroGradient="linear-gradient(135deg, #1a237e 0%, #283593 45%, #3949ab 100%)"
    icon={<GavelIcon sx={{ fontSize: 'inherit' }} />}
    tags={['ETPs', 'CriticalThinking', 'Epistemics', 'AIinEducation', 'PolicyDebate', 'KnowledgeGraphs', 'CHISG', 'Misinformation', 'EdTech']}
  >
    <Typography variant="body1" component="p">
      Real-world disagreement is rarely decided by the strength of an argument. It's decided by who triggers whose threat response first.
    </Typography>

    <Typography variant="body1" component="p">
      That's not cynicism. It's the twelve spectra, doing what they do. A claim arrives, Threat Response fires before Integrity Logic gets a vote, and the position is set before the evidence is even read. Politics runs on personality and triggered reaction far more than it runs on argument. If that's the actual mechanism, then the fix isn't a better argument. It's a process that doesn't let triggered reaction substitute for evidence in the first place.
    </Typography>

    <Typography variant="body1" component="p">
      That's what this post is — not a spectrum, but something built <em>using</em> the spectra. The first applied piece in this series.
    </Typography>

    <Typography variant="h2" component="h2">The Experiment</Typography>

    <Typography variant="body1" component="p">
      Take a real contested policy question — climate policy was the test case. Build the best, most honestly sourced version of each side: one arguing urgent mitigation, one arguing a gradual, cost-disciplined transition. Real citations. No straw men. Each side gets a full rebuttal of the other's strongest point, not just its weakest.
    </Typography>

    <Typography variant="body1" component="p">
      Then audit both sides with the same instrument.
    </Typography>

    <Typography variant="h2" component="h2">What the Governor Checks</Typography>

    <Typography variant="body1" component="p">
      The Governor isn't a debater and isn't a judge. It's an auditor. For every claim, on both sides, it checks four things:
    </Typography>

    <Box component="ul" sx={{ pl: 3, mb: 2, '& li': { mb: 1, color: 'text.secondary', lineHeight: 1.7 } }}>
      <li><strong>Logical consistency</strong> — does the conclusion actually follow from what was cited?</li>
      <li><strong>Evidence integrity</strong> — is the source saying what the claim says it's saying?</li>
      <li><strong>Source credibility</strong> — what kind of source is this, and what's its track record?</li>
      <li><strong>Rebuttal quality</strong> — did this side engage the other's actual best point, or a weaker stand-in?</li>
    </Box>

    <Typography variant="body1" component="p">
      One rule sits above all four, and it's the rule that makes the whole thing usable rather than just another voice in the argument:
    </Typography>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        No winner is declared. Ever.
      </Typography>
    </Box>

    <Typography variant="body1" component="p">
      Declaring a winner is exactly the move that turns an audit into another combatant — something to be argued with, not learned from. The Governor's job ends at showing what's supported, what's contested, and what's missing. What anyone does with that is theirs to decide.
    </Typography>

    <Typography variant="h2" component="h2">Tagging the Trigger, Not Diagnosing the Person</Typography>

    <Typography variant="body1" component="p">
      Every opening argument in the test case carried a field reading roughly: "many people feel <em>this</em> because of <em>that</em>." For the urgency side: Threat Response and Care Response. For the gradual-transition side: Risk Tolerance, plus a fairness concern.
    </Typography>

    <Typography variant="body1" component="p">
      Notice what that field is doing and what it isn't. It isn't saying the people holding the urgency position are running hot Threat Response as a trait, or that the gradualists are dispositionally risk-averse. It's naming which spectrum a <em>position</em> tends to recruit — separate from any individual's resting configuration. The same person can write the urgency case on Monday and feel the fairness concern on Tuesday. The tag describes the argument's emotional center of gravity, not a verdict on the arguer.
    </Typography>

    <Typography variant="body1" component="p">
      That distinction matters more than it looks. The moment a "many people feel" field becomes a fixed diagnosis of who's saying it, the tool stops being an audit and starts being a different kind of weapon — the same problem it was built to get away from.
    </Typography>

    <Typography variant="h2" component="h2">The Gap Map</Typography>

    <Typography variant="body1" component="p">
      The output isn't a score. It's a map, sorted into five categories per claim: Supported, Contested, Unsupported-as-stated, Contradicted, and Unresolved. Reading the map for the climate case, the IEA jobs claim landed in "unsupported-as-stated" — not because it was a lie, but because the report said something narrower and more conditional than the talking point built on top of it. That's the most common failure mode in real disagreement: not fabrication, but compression. A real, sourced finding gets summarized down to a punchier claim that the source doesn't quite carry.
    </Typography>

    <Typography variant="body1" component="p">
      Showing <em>why</em> the compressed version is easy to believe — what's true enough about the underlying finding to make the overstatement feel safe — does more work than just flagging it false. It explains the misreading instead of just penalizing it.
    </Typography>

    <Typography variant="h2" component="h2">Two Different Disagreements, Wearing the Same Clothes</Typography>

    <Typography variant="body1" component="p">
      The gap map surfaces something most debates never separate: two completely different kinds of disagreement, dressed identically as "we disagree about X."
    </Typography>

    <Typography variant="body1" component="p">
      <strong>Type one</strong> is factually resolvable but socially resisted. The evidence exists. The disagreement persists anyway, usually because accepting it costs someone status, money, or identity. This is a transparency problem — chip away at it with sourcing, with repetition, with the Governor's plain "this claim, in this source, says this much and no more."
    </Typography>

    <Typography variant="body1" component="p">
      <strong>Type two</strong> is a genuine values tension. Both sides see the same evidence and weigh it differently because they're optimizing for different things — speed versus stability, certainty versus cost. No amount of sourcing resolves this, because it was never a sourcing problem. This needs a mutual-legitimacy settlement: each side gets a real, falsifiable, time-bound commitment, not a forced merger of values that don't actually agree.
    </Typography>

    <Typography variant="body1" component="p">
      <strong>Old frame:</strong> "We just disagree — must be a values thing, let's agree to disagree."
    </Typography>
    <Typography variant="body1" component="p">
      <strong>New frame:</strong> "Is this resistance to evidence, or a genuine difference in what we're each protecting? Those need completely different tools."
    </Typography>

    <Typography variant="body1" component="p">
      Most public disagreement collapses these two together and then applies the wrong tool to both — endless re-litigating of settled facts (a type-one problem treated as type two), or false demands for consensus on a real values trade-off (a type-two problem treated as type one).
    </Typography>

    <Typography variant="h2" component="h2">What's Actually New Here</Typography>

    <Typography variant="body1" component="p">
      Worth being honest about what this is and isn't, because the honest answer is more useful than the flattering one.
    </Typography>

    <Typography variant="body1" component="p">
      The diagnostic layer — naming which spectrum a position recruits — is a direct reuse of the ETP framework this whole series has been building. Not a new idea wearing a new name; the same sliders, applied to public argument instead of personal regulation.
    </Typography>

    <Typography variant="body1" component="p">
      The settlement mechanics — falsifiable commitments, sunset clauses, the type-one/type-two split — are standard negotiation and epistemics theory (real-options reasoning, pre-registration, interest-based bargaining), competently assembled, not invented here.
    </Typography>

    <Typography variant="body1" component="p">
      The Governor itself — an evidence-and-logic auditor that explicitly refuses to adjudicate values, paired with a diagnostic layer that explains <em>why</em> each side believes what it believes — is the actual new piece. Neither half does much alone. Together they do something neither the ETP framework nor standard negotiation theory does on its own: explain the emotional pull of a position and check its factual scaffolding in the same pass, without letting either contaminate the other.
    </Typography>

    <Typography variant="h2" component="h2">Where CHISG Comes In — and Where It Stops</Typography>

    <Typography variant="body1" component="p">
      The Governor's "evidence integrity" and "source credibility" checks are, right now, judgment calls. CHISG — the knowledge-graph engine behind this site's skills mapping — has a formal version of exactly that: a computed trust score built from independent sources, independent logical routes to the same conclusion, and an echo-chamber filter that stops ten citations of one paper from counting as ten confirmations.
    </Typography>

    <Typography variant="body1" component="p">
      If that score existed for a claim already in the graph, the Governor wouldn't need to eyeball a source's credibility — it would read a number, with the method behind the number fully exposed.
    </Typography>

    <Typography variant="body1" component="p">
      Here's the boundary that has to hold, though. CHISG's trust scoring — including the part that spots when the same structural pattern recurs across unrelated domains — was built for mechanistic, settled-science claims: differential expansion causing curvature, recurring identically in a bimetallic strip and a stomatal cell. That recurrence really is evidence of an underlying law.
    </Typography>

    <Typography variant="body1" component="p">
      A claim like "this policy costs jobs" doesn't live in that regime. Two well-sourced experts can model the same economy and land in different places — that's not a gap in the graph, it's a genuine type-two disagreement wearing economic data as clothing. Treating the absence of a cross-domain structural match as evidence the claim is false would borrow CHISG's rigor without earning it — and would quietly reward whichever side cites more papers, which is precisely the failure mode the no-winner rule exists to prevent.
    </Typography>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        Use the trust score where a claim is actually checkable. Don't let it referee where the disagreement was never about facts.
      </Typography>
    </Box>

    <Typography variant="h2" component="h2">The Pattern, Applied</Typography>

    <Typography variant="body1" component="p">
      Twelve spectra, one diagnostic layer. The Governor is the first place this series has pointed the layer outward — not at a single nervous system trying to regulate itself, but at a public argument trying to find out what it's actually disagreeing about.
    </Typography>

    <Box component="blockquote">
      <Typography variant="body1" component="p">
        Not changing anyone's mind. Changing what's visible while they make it up.
      </Typography>
    </Box>
  </BlogPostLayout>
);

export default TheGovernor;
