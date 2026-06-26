import React, { useState, useRef, useEffect } from "react";
import { Container, Typography, Button, Paper, Box, Grid, IconButton, Tooltip, Dialog, DialogTitle, DialogContent, DialogContentText, DialogActions, Chip, Stack, Divider, alpha } from "@mui/material";
import UndoIcon from "@mui/icons-material/Undo";
import RedoIcon from "@mui/icons-material/Redo";
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined';
import RestartAltIcon from '@mui/icons-material/RestartAlt';
import SpectrumSlider from "../components/ETP/SpectrumSlider";
import { post } from "../api"; 

interface ETPSpectrum {
    id: string;
    label: string;
    description: string;
    negativeEffect: string;
    positiveEffect: string;
    conflictPattern: string;
    leftLabel: string;
    rightLabel: string;
}

// 12 Core biological spectra - innate, morally neutral, orthogonal
const spectra: ETPSpectrum[] = [
    {
        id: "social_gravity",
        label: "Social Gravity",
        description: "How social interaction affects your energy voltage.",
        negativeEffect: "Social interaction DRAINS voltage",
        positiveEffect: "Social interaction CHARGES voltage",
        conflictPattern: "Independent feels overwhelmed, Cohesive feels lonely",
        leftLabel: "Independent",
        rightLabel: "Cohesive"
    },
    {
        id: "energy_directionality",
        label: "Energy Directionality",
        description: "Where your energy naturally flows - inward or outward.",
        negativeEffect: "External stimulation OVERLOADS system",
        positiveEffect: "Internal energy NEEDS external outlet",
        conflictPattern: "Introverts drain, Extroverts starve",
        leftLabel: "Inward",
        rightLabel: "Outward"
    },
    {
        id: "voltage_sensitivity",
        label: "Voltage Sensitivity",
        description: "How much emotional insulation you need.",
        negativeEffect: "Needs thick emotional insulation",
        positiveEffect: "Thrives on emotional current",
        conflictPattern: "Insulated seems numb, Conductive seems dramatic",
        leftLabel: "Insulated",
        rightLabel: "Conductive"
    },
    {
        id: "threat_response",
        label: "Threat Response",
        description: "How you respond when threats are detected.",
        negativeEffect: "Threats trigger AVOIDANCE voltage (freeze/flee)",
        positiveEffect: "Threats trigger CONFRONTATION voltage (fight)",
        conflictPattern: "Passive seen as cowardly, Aggressive seen as dangerous",
        leftLabel: "Passive",
        rightLabel: "Aggressive"
    },
    {
        id: "care_response",
        label: "Care Response",
        description: "How you respond to vulnerability in others.",
        negativeEffect: "Vulnerability triggers DISTANCE voltage (self-protection)",
        positiveEffect: "Vulnerability triggers CARE voltage (other-protection)",
        conflictPattern: "Detached seen as cold, Nurturing seen as smothering",
        leftLabel: "Detached",
        rightLabel: "Nurturing"
    },
    {
        id: "risk_tolerance",
        label: "Risk Tolerance",
        description: "How risk affects your emotional state.",
        negativeEffect: "Risk creates ANXIETY voltage",
        positiveEffect: "Risk creates EXCITEMENT voltage",
        conflictPattern: "Averse holds back, Seeking pushes forward",
        leftLabel: "Averse",
        rightLabel: "Seeking"
    },
    {
        id: "integrity_logic",
        label: "Integrity Logic",
        description: "How you process moral and ethical decisions.",
        negativeEffect: "Absolutes create CONSTRICTION voltage",
        positiveEffect: "Relativity creates CHAOS voltage",
        conflictPattern: "Relativistic seems unprincipled, Absolutist seems rigid",
        leftLabel: "Relativistic",
        rightLabel: "Absolutist"
    },
    {
        id: "mirror_neuron_tuning",
        label: "Mirror Neuron Tuning",
        description: "How much you absorb others' emotional states.",
        negativeEffect: "Others' emotions are DISTINCT signals",
        positiveEffect: "Others' emotions are SHARED experience",
        conflictPattern: "Selective seems uncaring, Absorbent seems overwhelmed",
        leftLabel: "Selective",
        rightLabel: "Absorbent"
    },
    {
        id: "orderliness",
        label: "Orderliness",
        description: "Your preference for predictable structure versus spontaneous flexibility.",
        negativeEffect: "Structure creates CONSTRICTION voltage (feels caged)",
        positiveEffect: "Disorder creates ANXIETY voltage (feels chaotic)",
        conflictPattern: "Flexible seems chaotic, Ordered seems controlling",
        leftLabel: "Flexible",
        rightLabel: "Ordered"
    },
    {
        id: "responsibility_threshold",
        label: "Responsibility Threshold",
        description: "How you distribute blame and ownership when things go wrong.",
        negativeEffect: "Blame is EXTERNALISED (self-protection)",
        positiveEffect: "Blame is INTERNALISED (self-punishment)",
        conflictPattern: "Deflecting seems irresponsible, Absorbing seems martyrish",
        leftLabel: "Deflecting",
        rightLabel: "Absorbing"
    },
    {
        id: "loss_sensitivity",
        label: "Loss Sensitivity",
        description: "How you respond to losing possessions, status, or relationships.",
        negativeEffect: "Loss creates NO voltage (detachment)",
        positiveEffect: "Loss creates CRISIS voltage (territorial defence)",
        conflictPattern: "Detached seems indifferent, Territorial seems possessive",
        leftLabel: "Detached",
        rightLabel: "Territorial"
    },
    {
        id: "libido",
        label: "Libido",
        description: "The intensity of raw drive energy — wanting, ambition, appetite, desire.",
        negativeEffect: "Drive energy is CONTAINED (patience, restraint)",
        positiveEffect: "Drive energy DEMANDS expression (pursuit, urgency)",
        conflictPattern: "Restrained seems passionless, Expressive seems impulsive",
        leftLabel: "Restrained",
        rightLabel: "Expressive"
    },
];

// ---------- Condition presets ----------
// Each preset maps spectrum IDs to typical slider positions for that condition.
// Values not listed default to 0 (neurotypical midpoint).

interface ConditionPreset {
    id: string;
    label: string;
    subtitle: string;
    color: string;
    description: string;
    values: Partial<Record<string, number>>;
}

const conditionPresets: ConditionPreset[] = [
    {
        id: "autism",
        label: "Selective Processor",
        subtitle: "Autism Spectrum",
        color: "#5c6bc0",
        description: "Selective mirror tuning, high orderliness preference, independent social gravity. Optimised for deep, focused processing in predictable environments.",
        values: {
            social_gravity: -2,
            voltage_sensitivity: 2,
            mirror_neuron_tuning: -2,
            orderliness: 2,
            integrity_logic: 2,
        },
    },
    {
        id: "dyslexia",
        label: "Routing Bottleneck",
        subtitle: "Dyslexia",
        color: "#26a69a",
        description: "Primarily a skills-graph cluster (phonological routing), not a strong spectral pattern. Frustration voltage spikes from repeated decoding failure. Risk of internalised blame.",
        values: {
            energy_directionality: -1,
            voltage_sensitivity: 1,
            responsibility_threshold: 1,
        },
    },
    {
        id: "adhd",
        label: "Unfiltered Circuit",
        subtitle: "ADHD",
        color: "#ef5350",
        description: "Circuit runs hot and outward. High drive energy, risk-seeking, structure-resisting. Inhibitory control is the bottleneck — the pause between impulse and action is very short.",
        values: {
            voltage_sensitivity: 2,
            energy_directionality: 2,
            risk_tolerance: 2,
            orderliness: -2,
            libido: 2,
        },
    },
    {
        id: "audhd",
        label: "Compound Configuration",
        subtitle: "AuDHD (Autism + ADHD)",
        color: "#ab47bc",
        description: "Competing demands: needs order AND resists structure. Needs predictability AND craves novelty. The internal conflict itself generates voltage. Often presents as 'contradictory' behaviour.",
        values: {
            social_gravity: -2,
            voltage_sensitivity: 2,
            mirror_neuron_tuning: -2,
            integrity_logic: 2,
            energy_directionality: 1,
            risk_tolerance: 1,
            orderliness: 0,  // the conflict: autism wants +2, ADHD wants -2
            libido: 2,
        },
    },
    {
        id: "anxiety",
        label: "Overloaded Sentinel",
        subtitle: "Generalised Anxiety",
        color: "#ffa726",
        description: "Threat detection runs at maximum gain. Risk is feared, loss is catastrophised, responsibility is absorbed. The safety-maintenance loop dominates — everything is scanned for danger.",
        values: {
            voltage_sensitivity: 2,
            threat_response: -2,
            risk_tolerance: -2,
            responsibility_threshold: 2,
            loss_sensitivity: 2,
        },
    },
    {
        id: "neurotypical",
        label: "Midline Configuration",
        subtitle: "Neurotypical Baseline",
        color: "#78909c",
        description: "All spectra near the midpoint. No single circuit dominates or is suppressed. This isn't 'normal' — it's one configuration among many.",
        values: {},
    },
];

const MAX_HISTORY = 50;

const ETPProfilePage: React.FC = () => {
    const initialValues: Record<string, number> = {};
    spectra.forEach(s => (initialValues[s.id] = 0));

    const [values, setValues] = useState<Record<string, number>>(initialValues);
    const [history, setHistory] = useState<Record<string, number>[]>([initialValues]);
    const [historyIndex, setHistoryIndex] = useState(0);
    const [loading, setLoading] = useState(false);
    const [activePreset, setActivePreset] = useState<string | null>(null);

    const applyPreset = (preset: ConditionPreset) => {
        const newValues: Record<string, number> = {};
        spectra.forEach(s => {
            newValues[s.id] = preset.values[s.id] ?? 0;
        });
        setValues(newValues);
        setActivePreset(preset.id);

        // Push to history
        let newHistory = history.slice(0, historyIndex + 1);
        newHistory.push(newValues);
        if (newHistory.length > MAX_HISTORY) newHistory = newHistory.slice(newHistory.length - MAX_HISTORY);
        setHistory(newHistory);
        setHistoryIndex(newHistory.length - 1);
    };

    const resetAll = () => {
        setValues(initialValues);
        setActivePreset(null);
        let newHistory = history.slice(0, historyIndex + 1);
        newHistory.push(initialValues);
        if (newHistory.length > MAX_HISTORY) newHistory = newHistory.slice(newHistory.length - MAX_HISTORY);
        setHistory(newHistory);
        setHistoryIndex(newHistory.length - 1);
    };

    const handleValueChange = (id: string, val: number) => {
        setValues(prev => ({ ...prev, [id]: val }));
        setActivePreset(null); // manual adjustment clears preset indicator
    };

    // Called when user finishes a slider interaction
    const handleCommit = (id: string, val: number) => {
        // Values should already reflect the committed change
        const snapshot = { ...values };
        const last = history[history.length - 1];
        if (JSON.stringify(snapshot) === JSON.stringify(last)) return; // no change

        // Truncate future history if we had undone
        let newHistory = history.slice(0, historyIndex + 1);
        newHistory.push(snapshot);
        if (newHistory.length > MAX_HISTORY) newHistory = newHistory.slice(newHistory.length - MAX_HISTORY);

        setHistory(newHistory);
        setHistoryIndex(newHistory.length - 1);
    };

    const undo = () => {
        if (historyIndex <= 0) return;
        const newIndex = historyIndex - 1;
        setHistoryIndex(newIndex);
        setValues(history[newIndex]);
    };

    const redo = () => {
        if (historyIndex >= history.length - 1) return;
        const newIndex = historyIndex + 1;
        setHistoryIndex(newIndex);
        setValues(history[newIndex]);
    };

    // Keyboard shortcuts: Cmd/Ctrl+Z and Cmd/Ctrl+Y (redo)
    useEffect(() => {
        const onKey = (e: KeyboardEvent) => {
            const meta = e.ctrlKey || e.metaKey;
            if (!meta) return;
            if (e.key.toLowerCase() === 'z') {
                e.preventDefault();
                undo();
            } else if (e.key.toLowerCase() === 'y' || (e.shiftKey && e.key.toLowerCase() === 'z')) {
                e.preventDefault();
                redo();
            }
        };
        window.addEventListener('keydown', onKey);
        return () => window.removeEventListener('keydown', onKey);
    }, [history, historyIndex]);

    const handleSubmit = async () => {
        setLoading(true);
        try {
            console.log('Submitting profile:', values);
            const response = await post('/api/actionplan', { values });
            console.log('API Response:', response);
            if (response.status === 'success') {
                alert(`Action Plan Generated!\n\n${response.actionPlan.join('\n')}`);
            } else {
                alert('Failed to generate action plan.');
            }
        } catch (error) {
            console.error(error);
            alert('Error submitting profile. See console.');
        } finally {
            setLoading(false);
        }
    };

    const [profileHelpOpen, setProfileHelpOpen] = useState(false);

    return (
        <Container maxWidth="lg" sx={{ mt: 3, mb: 6 }}>
            <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 1 }}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <Typography variant="h4">ETP Personality Profile</Typography>
                    <Tooltip title="About this profile">
                        <IconButton size="small" aria-label="Profile help" onClick={() => setProfileHelpOpen(true)}>
                            <InfoOutlinedIcon fontSize="small" />
                        </IconButton>
                    </Tooltip>
                </Box>

                <Box>
                    <Tooltip title="Undo (Cmd/Ctrl+Z)">
                        <span>
                            <IconButton onClick={undo} disabled={historyIndex <= 0} aria-label="Undo">
                                <UndoIcon />
                            </IconButton>
                        </span>
                    </Tooltip>
                    <Tooltip title="Redo (Cmd/Ctrl+Y)">
                        <span>
                            <IconButton onClick={redo} disabled={historyIndex >= history.length - 1} aria-label="Redo">
                                <RedoIcon />
                            </IconButton>
                        </span>
                    </Tooltip>
                </Box>
            </Box>

            {/* Profile help dialog */}
            <Dialog open={profileHelpOpen} onClose={() => setProfileHelpOpen(false)} aria-labelledby="profile-help-title">
                <DialogTitle id="profile-help-title">About Emotional Trigger Points (ETP)</DialogTitle>
                <DialogContent>
                    <DialogContentText sx={{ mb: 2 }}>
                        ETP maps 12 core biological spectra that describe how your nervous system responds to different situations. These are innate, morally neutral, and independent of each other.
                    </DialogContentText>
                    <DialogContentText sx={{ mb: 2 }}>
                        Each spectrum has a natural "comfort zone" - neither end is better than the other. The goal is to expand your range of motion on each spectrum through targeted skills practice.
                    </DialogContentText>
                    <DialogContentText sx={{ mb: 1 }}>
                        Accelerator skills like <strong>executive function</strong> and <strong>self-regulation</strong> affect your capacity to deliberately move sliders across all spectra.
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setProfileHelpOpen(false)} color="primary">Close</Button>
                </DialogActions>
            </Dialog>

            {/* Condition presets */}
            <Paper sx={{ width: '100%', mb: 3, p: 2 }}>
                <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 1 }}>
                    <Typography variant="h6">
                        Condition Profiles
                    </Typography>
                    <Tooltip title="Reset all sliders to zero">
                        <IconButton size="small" onClick={resetAll}>
                            <RestartAltIcon fontSize="small" />
                        </IconButton>
                    </Tooltip>
                </Box>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                    Select a condition to see its typical ETP configuration. These are recognisable clusters — not diagnoses.
                    Each child's actual profile will vary. Adjust any slider after loading a preset to explore.
                </Typography>
                <Stack direction="row" spacing={1} sx={{ flexWrap: 'wrap', gap: 1 }}>
                    {conditionPresets.map((preset) => (
                        <Chip
                            key={preset.id}
                            label={preset.label}
                            onClick={() => applyPreset(preset)}
                            variant={activePreset === preset.id ? "filled" : "outlined"}
                            sx={{
                                fontWeight: activePreset === preset.id ? 700 : 400,
                                borderColor: preset.color,
                                color: activePreset === preset.id ? '#fff' : preset.color,
                                bgcolor: activePreset === preset.id ? preset.color : 'transparent',
                                '&:hover': {
                                    bgcolor: activePreset === preset.id ? preset.color : alpha(preset.color, 0.08),
                                },
                            }}
                        />
                    ))}
                </Stack>

                {/* Active preset description */}
                {activePreset && (() => {
                    const preset = conditionPresets.find(p => p.id === activePreset);
                    if (!preset) return null;
                    return (
                        <Box sx={{
                            mt: 2, p: 2, borderRadius: 2,
                            bgcolor: alpha(preset.color, 0.06),
                            borderLeft: `4px solid ${preset.color}`,
                        }}>
                            <Typography variant="subtitle2" sx={{ color: preset.color, mb: 0.5 }}>
                                {preset.subtitle}
                            </Typography>
                            <Typography variant="body2" color="text.secondary">
                                {preset.description}
                            </Typography>
                        </Box>
                    );
                })()}
            </Paper>

            {/* 12 Core Spectra */}
            <Paper sx={{ width: '100%', mb: 4, p: 2 }}>
                <Typography variant="h6" sx={{ mb: 2 }}>12 Core Spectra</Typography>
                <Grid container spacing={2}>
                    {spectra.map((s) => (
                        <Grid key={s.id} item xs={12} sm={6} md={4}>
                            <Box sx={{ p: 1 }}>
                                <SpectrumSlider
                                    id={s.id}
                                    title={s.label}
                                    description={s.description}
                                    value={values[s.id]}
                                    onChange={(val) => handleValueChange(s.id, val)}
                                    onCommit={(id, val) => handleCommit(id, val)}
                                    leftLabel={s.leftLabel}
                                    rightLabel={s.rightLabel}
                                />
                            </Box>
                        </Grid>
                    ))}
                </Grid>
            </Paper>

            <Button 
                variant="contained" 
                color="primary" 
                size="large" 
                fullWidth 
                onClick={handleSubmit} 
                disabled={loading}
            >
                {loading ? 'Generating Action Plan...' : 'Generate Learning & Action Plan'}
            </Button>
        </Container>
    );
};

export default ETPProfilePage;
