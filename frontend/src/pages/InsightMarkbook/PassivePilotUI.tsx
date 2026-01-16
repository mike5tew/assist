import React, { useState } from 'react';
import { Box, Paper, Typography, Grid, Slider, Modal, Button } from '@mui/material';

interface Student {
  id: string;
  name: string;
  position: { x: number; y: number };
  profile: { [key: string]: number }; // 15 Spectrums (-5 to 5)
}

const students: Student[] = [
  { id: '1', name: 'John Doe', position: { x: 0, y: 0 }, profile: { 'resilience': 0, 'empathy': 2 } },
  { id: '2', name: 'Jane Smith', position: { x: 1, y: 0 }, profile: { 'resilience': 1, 'empathy': -1 } },
  { id: '3', name: 'Alex Knight', position: { x: 2, y: 0 }, profile: { 'resilience': -2, 'empathy': 0 } },
  { id: '4', name: 'Sam Rivers', position: { x: 3, y: 0 }, profile: { 'resilience': 0, 'empathy': 1 } },
];

const BEHAVIOR_MARKERS = [
  { id: 'pos1', label: 'Positive Marker 1', icon: '⭐', color: '#4caf50' },
  { id: 'pos2', label: 'Positive Marker 2', icon: '🌟', color: '#81c784' },
  { id: 'neg1', label: 'Negative Marker 1', icon: '⚠️', color: '#f44336' },
  { id: 'neg2', label: 'Negative Marker 2', icon: '🛑', color: '#e57373' },
];

const SPECTRUMS = [
  'Shy — Outgoing', 'Guilt — Remorseless', 'Private — Expressive',
  'Aggressive — Passive', 'Empathic — Detached', 'Generous — Self-Int',
  'Bored — Enthusiastic', 'Patient — Impatient', 'Fragile — Resilient',
  'Risk-Averse — Seeking', 'Independent — Dependent', 'Conformist — Rebel',
  'Literal — Adaptable', 'Modest — Status', 'Truthful — Strategic'
];

const PassivePilotUI: React.FC = () => {
  const [lastAction, setLastAction] = useState<string>('System Ready');
  const [focusedStudentId, setFocusedStudentId] = useState<string | null>('1');
  const [showProfileSliders, setShowProfileSliders] = useState(false);
  const [temporalLog, setTemporalLog] = useState<any[]>([]);

  const logBehavior = (studentId: string, behaviorId: string) => {
    const behavior = BEHAVIOR_MARKERS.find(b => b.id === behaviorId);
    const student = students.find(s => s.id === studentId);
    const newEntry = {
      timestamp: new Date().toLocaleTimeString(),
      studentId,
      behaviorId,
      label: behavior?.label
    };
    setTemporalLog([newEntry, ...temporalLog]);
    setLastAction(`${student?.name}: ${behavior?.label} logged.`);
  };

  const handleSliderChange = (spectrum: string, value: number) => {
    // Logic for the "(Downwards) Rule" would be triggered here if value < current
    console.log(`Adjusting ${spectrum} to ${value}`);
  };

  return (
    <Box sx={{ p: 3, height: '100vh', display: 'flex', flexDirection: 'column' }}>
      <Typography variant="h4" gutterBottom>
        Passive Pilot: Behavioral Sampling
      </Typography>

      <Grid container spacing={3} sx={{ flexGrow: 1 }}>
        {/* Left Panel: Seating Plan */}
        <Grid item xs={8}>
          <Paper sx={{ p: 2, bgcolor: '#f5f5f5', height: '100%' }}>
            <Box sx={{ 
              display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 2 
            }}>
              {students.map(student => (
                <Paper 
                  key={student.id}
                  onClick={() => setFocusedStudentId(student.id)}
                  sx={{ 
                    height: '120px', 
                    display: 'flex', 
                    flexDirection: 'column',
                    alignItems: 'center', 
                    justifyContent: 'center',
                    cursor: 'pointer',
                    border: focusedStudentId === student.id ? '3px solid #1976d2' : '1px solid #ddd',
                    bgcolor: focusedStudentId === student.id ? '#e3f2fd' : 'white',
                    transition: 'all 0.2s'
                  }}
                >
                  <Typography variant="subtitle1" sx={{ fontWeight: 'bold' }}>
                    {student.name}
                  </Typography>
                  <Box sx={{ mt: 1, display: 'flex', gap: 0.5 }}>
                    {/* Mini status indicators could go here */}
                  </Box>
                </Paper>
              ))}
            </Box>
          </Paper>
        </Grid>

        {/* Right Panel: Focused Controls */}
        <Grid item xs={4}>
          <Paper sx={{ p: 3, height: '100%', display: 'flex', flexDirection: 'column' }}>
            {focusedStudentId ? (
              <>
                <Typography variant="h5" color="primary">
                  {students.find(s => s.id === focusedStudentId)?.name}
                </Typography>
                
                <Box sx={{ mt: 4 }}>
                  <Typography variant="overline" sx={{ letterSpacing: 1.5 }}>
                    Temporal Markers (Log Activity)
                  </Typography>
                  <Grid container spacing={2} sx={{ mt: 1 }}>
                    {BEHAVIOR_MARKERS.map(marker => (
                      <Grid item xs={6} key={marker.id}>
                        <Paper 
                          onClick={() => logBehavior(focusedStudentId, marker.id)}
                          sx={{ 
                            p: 2, textAlign: 'center', cursor: 'pointer',
                            transition: 'transform 0.1s',
                            '&:active': { transform: 'scale(0.95)' },
                            bgcolor: '#fafafa', 
                            border: `2px solid ${marker.color}22`,
                            '&:hover': { bgcolor: `${marker.color}11` }
                          }}
                        >
                          <Typography variant="h4" sx={{ color: marker.color }}>{marker.icon}</Typography>
                          <Typography variant="caption" sx={{ display: 'block', fontWeight: 'bold' }}>
                            {marker.label}
                          </Typography>
                        </Paper>
                      </Grid>
                    ))}
                  </Grid>
                </Box>

                <Box sx={{ mt: 4, flexGrow: 1 }}>
                  <Typography variant="overline" sx={{ letterSpacing: 1.5 }}>
                    Session Log
                  </Typography>
                  <Box sx={{ mt: 1, maxHeight: '150px', overflowY: 'auto' }}>
                    {temporalLog
                      .filter(log => log.studentId === focusedStudentId)
                      .map((log, i) => (
                        <Typography key={i} variant="caption" sx={{ display: 'block', borderBottom: '1px solid #f0f0f0', py: 0.5 }}>
                          [{log.timestamp}] {log.label}
                        </Typography>
                      ))}
                  </Box>
                </Box>

                <Box sx={{ mt: 'auto' }}>
                  <Paper 
                    onClick={() => setShowProfileSliders(true)}
                    sx={{ p: 2, textAlign: 'center', bgcolor: '#1976d2', color: 'white', cursor: 'pointer' }}
                  >
                    Adjust HumanOS Profile (Sliders)
                  </Paper>
                </Box>
              </>
            ) : (
              <Typography color="textSecondary">Select a student from the plan.</Typography>
            )}
          </Paper>
        </Grid>
      </Grid>
      
      {/* Modal for Profile Sliders */}
      <Modal open={showProfileSliders} onClose={() => setShowProfileSliders(false)}>
        <Box sx={{ 
          position: 'absolute', top: '50%', left: '50%', transform: 'translate(-50%, -50%)',
          width: 800, bgcolor: 'background.paper', p: 4, borderRadius: 2, boxShadow: 24,
          maxHeight: '90vh', overflowY: 'auto'
        }}>
          <Typography variant="h5" sx={{ mb: 3 }}>
            HumanOS Profile: {students.find(s => s.id === focusedStudentId)?.name}
          </Typography>
          <Grid container spacing={3}>
            {SPECTRUMS.map(spec => (
              <Grid item xs={6} key={spec}>
                <Typography variant="caption">{spec}</Typography>
                <Slider 
                  defaultValue={0} 
                  min={-5} max={5} step={1} 
                  marks={[{value: 0, label: '0'}]}
                  valueLabelDisplay="auto"
                  onChangeCommitted={(e, v) => handleSliderChange(spec, v as number)}
                />
              </Grid>
            ))}
          </Grid>
          <Box sx={{ mt: 4, textAlign: 'right' }}>
            <Button variant="contained" onClick={() => setShowProfileSliders(false)}>Save & Close</Button>
          </Box>
        </Box>
      </Modal>

      {/* Footer / Status Bar */}
      <Box sx={{ mt: 2, p: 1, bgcolor: '#333', color: 'white', borderRadius: 1 }}>
        <Typography variant="caption">Last Action: {lastAction}</Typography>
      </Box>
    </Box>
  );
};

export default PassivePilotUI;
