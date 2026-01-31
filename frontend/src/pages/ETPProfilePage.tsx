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
    example?: string;
    category: "Social" | "Emotional" | "Authority";
    leftLabel: string;
    rightLabel: string;
}

const spectra: ETPSpectrum[] = [
    // Social
    { id: "social_gravity", label: "Social Gravity", description: "Tendency to align with group social pull versus independence.", example: "E.g., prefers group consensus over going solo on tasks.", category: "Social", leftLabel: "Independent", rightLabel: "Cohesive" },
    { id: "guilt_response", label: "Guilt Response", description: "How quickly guilt is internalized versus externalized in social interactions.", example: "E.g., feels guilty quickly after a perceived social misstep.", category: "Social", leftLabel: "Internalized", rightLabel: "Externalized" },
    { id: "emotional_transparency", label: "Emotional Transparency", description: "Openness in expressing emotional state to others.", example: "E.g., openly shares worry vs keeping feelings private.", category: "Social", leftLabel: "Opaque", rightLabel: "Transparent" },
    { id: "energy_directionality", label: "Energy Directionality", description: "Whether energy is directed inward (reflective) or outward (expressive).", example: "E.g., energizes others in group vs prefers quiet reflection.", category: "Social", leftLabel: "Inward (Intro)", rightLabel: "Outward (Extro)" },
    { id: "mirror_neuron_tuning", label: "Mirror Neuron Tuning", description: "Sensitivity to others' actions and emotions leading to mirroring behavior.", example: "E.g., unconsciously mirrors colleagues' body language.", category: "Social", leftLabel: "Selective", rightLabel: "Absorbent" },
    { id: "resource_allocation", label: "Resource Allocation", description: "Tendency to hoard versus share social and material resources.", example: "E.g., keeps learning materials private vs shares freely.", category: "Social", leftLabel: "Hoarding", rightLabel: "Sharing" },

    // Emotional
    { id: "voltage_sensitivity", label: "Voltage Sensitivity", description: "Emotional reactivity to small inputs—how easily one is affected.", example: "E.g., small criticisms either deeply unsettle or are shrugged off.", category: "Emotional", leftLabel: "Insulated", rightLabel: "Conductive" },
    { id: "impulse_gap", label: "Impulse Gap", description: "Ability to pause and reflect before acting versus reacting immediately.", example: "E.g., takes a breath before replying vs answers instantly.", category: "Emotional", leftLabel: "Reactive", rightLabel: "Reflective" },
    { id: "self_righting_speed", label: "Self-Righting Speed", description: "Speed at which emotional state returns to baseline after disturbance.", example: "E.g., recovers quickly after setbacks vs dwelling for days.", category: "Emotional", leftLabel: "Slow", rightLabel: "Rapid" },
    { id: "risk_tolerance", label: "Risk Tolerance", description: "Comfort with uncertainty and willingness to take risks.", example: "E.g., volunteers for ambitious projects vs prefers safe options.", category: "Emotional", leftLabel: "Averse", rightLabel: "Seeking" },
    { id: "anticipation_bias", label: "Anticipation Bias", description: "Lean towards optimistic or pessimistic anticipation of outcomes.", example: "E.g., expects success vs prepares for problems.", category: "Emotional", leftLabel: "Pessimistic", rightLabel: "Optimistic" },
    { id: "presence_sensitivity", label: "Presence Sensitivity", description: "Degree of attunement to the present moment and surroundings.", example: "E.g., notices subtle cues in conversation vs seems distracted.", category: "Emotional", leftLabel: "Detached", rightLabel: "Attuned" },

    // Authority
    { id: "agency_threshold", label: "Agency Threshold", description: "Threshold for taking action versus waiting for direction.", example: "E.g., starts tasks proactively vs waits for instructions.", category: "Authority", leftLabel: "Passive", rightLabel: "Active" },
    { id: "authority_response", label: "Authority Response", description: "Disposition towards compliance or challenging authority.", example: "E.g., follows instructions strictly vs questions policies.", category: "Authority", leftLabel: "Compliant", rightLabel: "Challenging" },
    { id: "ambiguity_tolerance", label: "Ambiguity Tolerance", description: "Comfort operating under unclear or incomplete information.", example: "E.g., works well with open-ended tasks vs needs clear specs.", category: "Authority", leftLabel: "Rigid", rightLabel: "Fluid" },
    { id: "status_sensitivity", label: "Status Sensitivity", description: "Degree to which social status affects decisions and behavior.", example: "E.g., defers to seniority vs treats everyone equally.", category: "Authority", leftLabel: "Indifferent", rightLabel: "Concerned" },
    { id: "integrity_logic", label: "Integrity Logic", description: "Preference for relativistic versus absolutist moral reasoning.", example: "E.g., balances context in decisions vs applies strict rules.", category: "Authority", leftLabel: "Relativistic", rightLabel: "Absolutist" },
];

const MAX_HISTORY = 50;

const ETPProfilePage: React.FC = () => {
    const initialValues: Record<string, number> = {};
    spectra.forEach(s => (initialValues[s.id] = 0));

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
                        ETP stands for Emotional Trigger Points. This system maps learner approaches across Social, Emotional, and Authority domains.
                    </DialogContentText>
                    <DialogContentText sx={{ mb: 1 }}>
                        This page displays the logic and approaches of the system using 17 modulation spectra. The focus is on visualising behavioural balances and identifying qualitative patterns, rather than generating numerical scores.
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setProfileHelpOpen(false)} color="primary">Close</Button>
                </DialogActions>
            </Dialog>

            <Paper sx={{ width: '100%', mb: 4, p: 2 }}>
                <Grid container spacing={2}>
                    {spectra.map((s) => (
                        <Grid key={s.id} item xs={12} sm={6} md={4}>
                            <Box sx={{ p: 1 }}>
                                <SpectrumSlider
                                    id={s.id}
                                    title={s.label}
                                    description={s.description}
                                    example={s.example}
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
