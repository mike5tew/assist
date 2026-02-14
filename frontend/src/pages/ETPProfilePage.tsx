import React, { useState, useRef, useEffect } from "react";
import { Container, Typography, Button, Paper, Box, Grid, IconButton, Tooltip, Dialog, DialogTitle, DialogContent, DialogContentText, DialogActions } from "@mui/material";
import UndoIcon from "@mui/icons-material/Undo";
import RedoIcon from "@mui/icons-material/Redo";
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined';
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

interface GlobalModerator {
    id: string;
    label: string;
    description: string;
    lowEffect: string;
    highEffect: string;
}

// 8 Core biological spectra - innate, morally neutral, orthogonal
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
];

// Global moderators affect ALL spectra
const globalModerators: GlobalModerator[] = [
    {
        id: "pilot_strength",
        label: "Pilot Strength",
        description: "Executive function capacity - the hand on ALL sliders.",
        lowEffect: "Sliders move reactively; triggers control positioning",
        highEffect: "Sliders move deliberately; conscious choice of positioning"
    },
    {
        id: "current_load",
        label: "Current Load",
        description: "Stress/depletion level - narrows range of motion on all spectra.",
        lowEffect: "Full range of motion available on all spectra",
        highEffect: "Range contracts toward comfort zone; less flexibility"
    },
];

const MAX_HISTORY = 50;

const ETPProfilePage: React.FC = () => {
    const initialValues: Record<string, number> = {};
    spectra.forEach(s => (initialValues[s.id] = 0));
    globalModerators.forEach(m => (initialValues[m.id] = 0));

    const [values, setValues] = useState<Record<string, number>>(initialValues);
    const [history, setHistory] = useState<Record<string, number>[]>([initialValues]);
    const [historyIndex, setHistoryIndex] = useState(0);
    const [loading, setLoading] = useState(false);

    const handleValueChange = (id: string, val: number) => {
        setValues(prev => ({ ...prev, [id]: val }));
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
                        ETP maps 8 core biological spectra that describe how your nervous system responds to different situations. These are innate, morally neutral, and independent of each other.
                    </DialogContentText>
                    <DialogContentText sx={{ mb: 2 }}>
                        Each spectrum has a natural "comfort zone" - neither end is better than the other. The goal is to expand your range of motion on each spectrum through targeted skills practice.
                    </DialogContentText>
                    <DialogContentText sx={{ mb: 1 }}>
                        Two global moderators affect ALL spectra: <strong>Pilot Strength</strong> (your executive function capacity) and <strong>Current Load</strong> (stress level that narrows your range).
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setProfileHelpOpen(false)} color="primary">Close</Button>
                </DialogActions>
            </Dialog>

            {/* Global Moderators */}
            <Paper sx={{ width: '100%', mb: 3, p: 2, bgcolor: 'grey.50' }}>
                <Typography variant="h6" sx={{ mb: 2 }}>Global Moderators</Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                    These affect your capacity across ALL spectra
                </Typography>
                <Grid container spacing={2}>
                    {globalModerators.map((m) => (
                        <Grid key={m.id} item xs={12} sm={6}>
                            <Box sx={{ p: 1 }}>
                                <SpectrumSlider
                                    id={m.id}
                                    title={m.label}
                                    description={m.description}
                                    value={values[m.id] || 0}
                                    onChange={(val) => handleValueChange(m.id, val)}
                                    onCommit={(id, val) => handleCommit(id, val)}
                                    leftLabel="Low"
                                    rightLabel="High"
                                />
                            </Box>
                        </Grid>
                    ))}
                </Grid>
            </Paper>

            {/* 8 Core Spectra */}
            <Paper sx={{ width: '100%', mb: 4, p: 2 }}>
                <Typography variant="h6" sx={{ mb: 2 }}>8 Core Spectra</Typography>
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
