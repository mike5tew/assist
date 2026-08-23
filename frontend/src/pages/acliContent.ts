export interface ACLILesson {
  id: number;
  title: string;
  content: string;
}

export interface ACLIModule {
  id: number;
  name: string;
  lessons: ACLILesson[];
}

export const ACLI_GLOSSARY: Record<string, string> = {
  "V(D)J recombination": "Somatic gene-segment rearrangement that builds antigen receptor variable regions in developing B and T cells.",
  "RAG1/2": "Lymphoid recombinase complex that recognizes RSS and initiates DNA cleavage for receptor assembly.",
  "RSS": "Recombination signal sequence flanking V, D, and J segments; contains conserved heptamer/nonamer with 12 or 23 bp spacer.",
  "12/23 rule": "Recombination usually occurs between one 12-bp spacer RSS and one 23-bp spacer RSS to enforce correct segment joining.",
  "c-NHEJ": "Classical non-homologous end joining DNA repair pathway that rejoins programmed V(D)J DNA breaks.",
  "Artemis": "DNA endonuclease activated by DNA-PKcs that opens hairpin coding ends after RAG cleavage.",
  "DNA-PKcs": "Catalytic kinase in DNA-PK complex that coordinates NHEJ and activates Artemis.",
  "TdT": "Template-independent DNA polymerase that adds non-templated nucleotides to coding ends, increasing CDR3 diversity.",
  "P nucleotides": "Palindromic nucleotides generated during asymmetric hairpin opening at coding ends.",
  "CDR3": "Most variable part of BCR/TCR variable region, formed directly at V(D)J junctions and central to specificity.",
  "TREC": "T-cell receptor excision circle; marker of recent thymic output and used in SCID newborn screening.",
  "KREC": "Kappa-deleting recombination excision circle; marker of recent B-cell output.",
  "T-B-NK+ SCID": "Immunophenotype with absent T and B cells but preserved NK cells, classically seen in complete V(D)J pathway defects.",
  "Hypomorphic": "Partial loss-of-function variant that retains residual activity and creates attenuated or dysregulated phenotypes."
};

export const ACLI_MODULES: ACLIModule[] = [
  {
    id: 5001,
    name: "Track 1: Foundation Build",
    lessons: [
      {
        id: 500101,
        title: "Speedread Stage Map",
        content: `# Speedread: V(D)J Pathway in 9 Steps

## Core sequence
1. Locus accessibility in early lymphocyte development (G0/G1-biased).  
2. Segment architecture available: V, D, J (or V, J only depending on locus).  
3. RSS pairing and the [[12/23 rule|12/23 rule]].  
4. [[RAG1/2|RAG1/2]] cleavage at RSS-coding boundary.  
5. Hairpin coding ends + blunt signal ends are formed.  
6. [[Artemis|Artemis]] opening of coding hairpins (with [[DNA-PKcs|DNA-PKcs]] support).  
7. End processing: trimming, [[P nucleotides|P nucleotides]], [[TdT|TdT]] N addition.  
8. Joining by [[c-NHEJ|c-NHEJ]] (Ku, DNA-PKcs, XRCC4, XLF, Ligase IV).  
9. Productive receptor expression and developmental selection.

## One-line clinical link
Complete failure of steps 3-8 in both lineages gives [[T-B-NK+ SCID|T-B-NK+ SCID]] when NK development is otherwise preserved.

## Mnemonic (pathway order)
"Access Segments, Signal RAG, Hairpins Open, Trim + TdT, NHEJ Locks".
`
      },
      {
  id: 500102,
        title: "Exam Core: Mechanism",
        content: `# Exam Core: How V(D)J Recombination Works

## Locus and timing
- Rearrangement occurs before antigen exposure during defined developmental windows.
- B-cell order: IGH D-J then V-DJ, then light chain V-J.
- T-cell order: TCR beta (and gamma/delta) first, then TCR alpha.

## Segment usage by locus
- VDJ loci: IGH, TRB, TRD
- VJ loci: IGK, IGL, TRA, TRG

## Cleavage and repair logic
- [[RSS|RSS]] defines legal recombination boundaries.
- [[RAG1/2|RAG1/2]] generates programmed double-strand breaks.
- Coding ends are hairpins; signal ends are blunt.
- [[c-NHEJ|c-NHEJ]] performs final ligation after processing.

## High-yield distinction
- RAG makes the break.
- Artemis opens the hairpin.
- NHEJ machinery seals the ends.

## Mnemonic (repair proteins)
"Ku-PK-Artemis-X4-XLF-L4"  
(ordered memory aid: Ku70/80, DNA-PKcs, Artemis, XRCC4, XLF, Ligase IV)
`
      },
      {
  id: 500103,
        title: "Diversity: Combinatorial vs Junctional",
        content: `# Diversity: What Creates the Huge Repertoire

## Combinatorial diversity
- Choice of V, D, J segments.
- Heavy/light pairing in BCR.
- Alpha/beta pairing in TCR.

## Junctional diversity
- Imprecise joins at coding ends.
- Exonuclease trimming.
- [[TdT|TdT]] non-templated nucleotide addition.
- [[P nucleotides|P nucleotides]] from hairpin opening.
- Maximal impact in [[CDR3|CDR3]].

## Why this matters in viva
Even with limited germline segment counts, multiplicative diversity and CDR3 variability create very large theoretical sequence space.

## Mnemonic
"Choose, Pair, then Sculpt the Junction".
`
      }
    ]
  },
  {
    id: 5002,
    name: "Track 2: Comparison and Clinical Spectrum",
    lessons: [
      {
        id: 500201,
        title: "BCR vs TCR Rearrangement",
        content: `# Compare and Contrast: B-cell vs T-cell Rearrangement

## Shared machinery
- Both use [[RAG1/2|RAG1/2]], [[RSS|RSS]], and [[c-NHEJ|c-NHEJ]].

## Key differences
- B-cell receptor editing is important for tolerance.
- TCR beta allelic exclusion is particularly critical.
- Selection context differs: pre-BCR/BCR vs pre-TCR plus thymic selection.

## Exam-ready table points
- VDJ usage: Ig heavy vs TCR beta/delta.
- VJ usage: Ig kappa/lambda vs TCR alpha/gamma.
- Initial productive chain: mu heavy vs typically beta in alpha-beta T lineage.

## Mnemonic
"Same tools, different timetable and checkpoints".
`
      },
      {
        id: 500202,
        title: "Clinical Application: SCID Spectrum",
        content: `# Clinical Application: SCID Patterns from the V(D)J Axis

## Classical complete recombination failure
- Phenotype: [[T-B-NK+ SCID|T-B-NK+ SCID]].
- Typical category: complete RAG deficiency or key c-NHEJ defects.
- Newborn clue: absent [[TREC|TREC]].

## Why NK cells can be preserved
NK cells do not require clonotypic V(D)J receptor rearrangement for development.

## Hypomorphic variants
- Residual activity ([[Hypomorphic|hypomorphic]]) can produce atypical SCID, CID, Omenn-like disease, granulomatous inflammation, and immune dysregulation.
- Expect restricted/oligoclonal repertoires rather than total absence.

## Additional lab support beyond subsets
- Repertoire sequencing for clonal restriction and CDR3 skew.
- [[TREC|TREC]]/[[KREC|KREC]] style output assays.
- Functional RAG assays where available.

## Mnemonic (complete vs partial)
"None gives absence; some gives skew".
`
      }
    ]
  },
  {
    id: 5003,
    name: "Track 3: Repair Defects and Interpretation Drill",
    lessons: [
      {
        id: 500301,
        title: "DNA Repair Defects and Radiosensitivity",
        content: `# DNA Repair Defects in the V(D)J Pathway

## Artemis checkpoint
- [[Artemis|Artemis]] opens coding-end hairpins after RAG cleavage.
- Failure blocks progression to ligation.

## Why radiosensitivity differs
- Isolated RAG defects are lymphoid recombination-specific.
- NHEJ defects (for example Artemis-pathway) impair broader double-strand break repair and can show cellular radiosensitivity.

## Useful differential genes
- DCLRE1C (Artemis)
- PRKDC (DNA-PKcs)
- LIG4 (DNA ligase IV)
- NHEJ1 (XLF)
- XRCC4

## Exam discriminator
T-B-NK+ with marked radiosensitivity strongly points to DNA-repair/NHEJ category over isolated RAG category.

## Mnemonic
"RAG is the cutter; NHEJ is the body-repair crew".
`
      },
      {
        id: 500302,
        title: "Interpretation Grid Drill",
        content: `# Interpretation Grid Drill

## Pattern mapping
- Absent T, absent B, present NK, no radiosensitivity -> RAG category likely.
- Absent T, absent B, present NK, radiosensitivity -> NHEJ/DNA repair category likely.
- Reduced or oligoclonal T with inflammatory dysregulation signs -> hypomorphic recombination defect likely.
- Broad lymphopenia plus metabolic abnormalities -> metabolic lymphocyte survival defect category.

## Representative genes
- RAG category: RAG1 or RAG2
- NHEJ category: DCLRE1C, PRKDC, LIG4, NHEJ1
- Metabolic category example: ADA

## High-yield caveat
Same gene can produce different phenotypes depending on residual protein function.

## Mnemonic
"Same genotype family, different residual-function story".
`
      }
    ]
  }
];
