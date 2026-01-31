import React, { useState } from 'react';
import { Box, Slider, Typography, Grid, IconButton, Dialog, DialogTitle, DialogContent, DialogContentText, DialogActions, Button } from '@mui/material';
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined';

interface SpectrumSliderProps {
    id: string;
    title?: string;
    description?: string;
    example?: string;
    value: number;
    onChange: (value: number) => void;
    onCommit?: (id: string, value: number) => void;
    leftLabel: string;
    rightLabel: string;
}

function backgroundForValue(val: number) {
    switch (val) {
        case -2:
            return '#ff4d4f'; // red
        case -1:
            return '#ffd666'; // yellow
        case 0:
            return '#73d13d'; // green
        case 1:
            return '#36cfc9'; // cyan
        case 2:
            return '#9254de'; // purple
        default:
            return 'transparent';
    }
}

const SpectrumSlider: React.FC<SpectrumSliderProps> = ({ 
    id,
    title,
    description,
    example,
    value, 
    onChange, 
    onCommit,
    leftLabel, 
    rightLabel 
}) => {
    const bg = backgroundForValue(value ?? 0);
    const [open, setOpen] = useState(false);

    return (
        <Box sx={{ width: '100%', mb: 0.5, position: 'relative' }}>
            {/* Behaviour labels above the slider, at each end, and a help icon */}
            <Grid container spacing={0.5} alignItems="center" sx={{ mb: 0.5 }}>
                <Grid item xs={5}>
                    <Typography variant="caption" align="left" component="div">
                        {leftLabel}
                    </Typography>
                </Grid>
                <Grid item xs={2} sx={{ display: 'flex', justifyContent: 'center' }}>
                    <IconButton size="small" aria-label={`Info for ${title || ''}`} onClick={() => setOpen(true)}>
                        <InfoOutlinedIcon fontSize="small" />
                    </IconButton>
                </Grid>
                <Grid item xs={5}>
                    <Typography variant="caption" align="right" component="div">
                        {rightLabel}
                    </Typography>
                </Grid>
            </Grid>

            {/* Colored background container */}
            <Box sx={{ background: bg, borderRadius: 1, p: 0.5 }}>
                <Slider
                    value={value}
                    onChange={(_, val) => onChange(val as number)}
                    step={1}
                    min={-2}
                    max={2}
                    valueLabelDisplay="off"
                    onChangeCommitted={(_, val) => onCommit && onCommit(id, val as number)}
                />
            </Box>

            <Dialog open={open} onClose={() => setOpen(false)} aria-labelledby={`dialog-${id}-title`}>
                <DialogTitle id={`dialog-${id}-title`}>{title}</DialogTitle>
                <DialogContent>
                    <DialogContentText sx={{ mb: 1 }}>
                        {description}
                    </DialogContentText>
                    {example && (
                        <DialogContentText sx={{ fontStyle: 'italic' }}>
                            Example: {example}
                        </DialogContentText>
                    )}
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setOpen(false)} color="primary">Close</Button>
                </DialogActions>
            </Dialog>
        </Box>
    );
};

export default SpectrumSlider;
