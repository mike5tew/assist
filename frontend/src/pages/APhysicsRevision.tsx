import React, { useState, useEffect, useRef, useCallback } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkMath from 'remark-math';
import remarkGfm from 'remark-gfm';
import rehypeKatex from 'rehype-katex';
import 'katex/dist/katex.min.css';
import {
  Table, TableHead, TableBody, TableRow, TableCell, TableContainer, Paper,
} from '@mui/material';
import { APHY_MODULES, APhysLesson, APhysModule } from './aPhysicsContent';

const MODULE_COLORS: Record<number, string> = {
  4001: '#1565c0',
  4002: '#6a1b9a',
  4003: '#2e7d32',
  4004: '#e65100',
};

type Mode = 'read' | 'speedread' | 'audio' | 'flashcards';

// ── Flashcard helpers ────────────────────────────────────────────────────────

interface Flashcard {
  front: string;
  back: string;
}

function extractFlashcards(markdown: string): Flashcard[] {
  const lines = markdown.split('\n');
  const cards: Flashcard[] = [];
  let currentHeading = '';
  let bodyLines: string[] = [];

  const flush = () => {
    if (currentHeading && bodyLines.length > 0) {
      const back = bodyLines
        .join('\n')
        .trim()
        .replace(/^\s*[-*]\s+/gm, '• ')
        .trim();
      if (back.length > 10) cards.push({ front: currentHeading, back });
    }
  };

  for (const line of lines) {
    const h2 = line.match(/^## (.+)/);
    const h3 = line.match(/^### (.+)/);
    if (h2 || h3) {
      flush();
      currentHeading = (h2 || h3)![1];
      bodyLines = [];
    } else {
      bodyLines.push(line);
    }
  }
  flush();
  return cards;
}

// ── Strip markdown to plain text for TTS / speed-read ───────────────────────

function markdownToPlainText(md: string): string {
  return md
    .replace(/```[\s\S]*?```/g, '')
    .replace(/\$\$[\s\S]*?\$\$/g, 'equation')
    .replace(/\$[^$]+\$/g, 'equation')
    .replace(/^\|.*\|$/gm, '')     // table rows
    .replace(/^#+\s+/gm, '')       // headings
    .replace(/\*\*([^*]+)\*\*/g, '$1')
    .replace(/\*([^*]+)\*/g, '$1')
    .replace(/`([^`]+)`/g, '$1')
    .replace(/^\s*[-*]\s+/gm, '')
    .replace(/\[([^\]]+)\]\([^)]+\)/g, '$1')
    .replace(/\n{3,}/g, '\n\n')
    .trim();
}

// ── Main component ────────────────────────────────────────────────────────────

export default function APhysicsRevision() {
  const [selectedModule, setSelectedModule] = useState<APhysModule>(APHY_MODULES[0]);
  const [selectedLesson, setSelectedLesson] = useState<APhysLesson>(APHY_MODULES[0].lessons[0]);
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [mode, setMode] = useState<Mode>('read');

  // Speed-read state
  const [srWords, setSrWords] = useState<string[]>([]);
  const [srIndex, setSrIndex] = useState(0);
  const [srPlaying, setSrPlaying] = useState(false);
  const [srComplete, setSrComplete] = useState(false);
  const [srWpm, setSrWpm] = useState(300);
  const srTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  // Audio state
  const [audioPlaying, setAudioPlaying] = useState(false);
  const utteranceRef = useRef<SpeechSynthesisUtterance | null>(null);

  // Flashcard state
  const [flashcards, setFlashcards] = useState<Flashcard[]>([]);
  const [fcIndex, setFcIndex] = useState(0);
  const [fcFlipped, setFcFlipped] = useState(false);

  const accentColor = MODULE_COLORS[selectedModule.id] || '#1565c0';

  const handleLessonSelect = useCallback((mod: APhysModule, lesson: APhysLesson) => {
    // Stop any active modes
    if (audioPlaying) {
      window.speechSynthesis?.cancel();
      setAudioPlaying(false);
    }
    if (srTimerRef.current) clearTimeout(srTimerRef.current);
    setSrPlaying(false);
    setSrComplete(false);
    setSrIndex(0);

    setSelectedModule(mod);
    setSelectedLesson(lesson);
    setMode('read');
    if (window.innerWidth < 768) setSidebarOpen(false);
  }, [audioPlaying]);

  // Initialise speed-reader words when lesson/mode changes
  useEffect(() => {
    if (mode === 'speedread') {
      const text = markdownToPlainText(selectedLesson.content);
      const words = text.split(/\s+/).filter(w => w.length > 0);
      setSrWords(words);
      setSrIndex(0);
      setSrPlaying(false);
      setSrComplete(false);
    }
  }, [mode, selectedLesson]);

  // Speed-reader ticker
  useEffect(() => {
    if (!srPlaying || srIndex >= srWords.length) return;
    const ms = 60000 / srWpm;
    srTimerRef.current = setTimeout(() => {
      if (srIndex < srWords.length - 1) {
        setSrIndex(i => i + 1);
      } else {
        setSrPlaying(false);
        setSrComplete(true);
      }
    }, ms);
    return () => { if (srTimerRef.current) clearTimeout(srTimerRef.current); };
  }, [srPlaying, srIndex, srWpm, srWords.length]);

  // Initialise flashcards when mode changes
  useEffect(() => {
    if (mode === 'flashcards') {
      setFlashcards(extractFlashcards(selectedLesson.content));
      setFcIndex(0);
      setFcFlipped(false);
    }
  }, [mode, selectedLesson]);

  const currentIdx = selectedModule.lessons.indexOf(selectedLesson);
  const prevLesson = currentIdx > 0 ? selectedModule.lessons[currentIdx - 1] : null;
  const nextLesson = currentIdx < selectedModule.lessons.length - 1 ? selectedModule.lessons[currentIdx + 1] : null;

  // ── Audio handlers ──────────────────────────────────────────────────────────

  const handleAudioPlay = () => {
    const text = markdownToPlainText(selectedLesson.content);
    const utt = new SpeechSynthesisUtterance(text);
    utt.rate = 0.95;
    utt.onend = () => setAudioPlaying(false);
    utteranceRef.current = utt;
    window.speechSynthesis.speak(utt);
    setAudioPlaying(true);
  };

  const handleAudioStop = () => {
    window.speechSynthesis?.cancel();
    setAudioPlaying(false);
  };

  // ── Flashcard handlers ──────────────────────────────────────────────────────

  const fcCard = flashcards[fcIndex];

  // ── Markdown components ─────────────────────────────────────────────────────

  const mdComponents: any = {
    h1: ({ children }: any) => <h1 style={{ fontSize: 22, fontWeight: 700, color: accentColor, marginBottom: 12, marginTop: 0 }}>{children}</h1>,
    h2: ({ children }: any) => <h2 style={{ fontSize: 17, fontWeight: 600, color: '#333', marginTop: 24, marginBottom: 8, borderBottom: `2px solid ${accentColor}20`, paddingBottom: 4 }}>{children}</h2>,
    h3: ({ children }: any) => <h3 style={{ fontSize: 15, fontWeight: 600, color: '#444', marginTop: 16, marginBottom: 6 }}>{children}</h3>,
    strong: ({ children }: any) => <strong style={{ color: accentColor }}>{children}</strong>,
    blockquote: ({ children }: any) => (
      <blockquote style={{ borderLeft: `4px solid ${accentColor}`, margin: '12px 0', padding: '8px 16px', background: `${accentColor}08`, borderRadius: '0 4px 4px 0', color: '#555' }}>
        {children}
      </blockquote>
    ),
    code: ({ inline, children }: any) => inline
      ? <code style={{ background: '#f0f0f0', padding: '1px 5px', borderRadius: 3, fontSize: 13, fontFamily: 'monospace' }}>{children}</code>
      : <pre style={{ background: '#1a1a2e', color: '#e8e8e8', padding: '12px 16px', borderRadius: 6, overflowX: 'auto', fontSize: 13 }}><code>{children}</code></pre>,
    table: ({ children }: any) => (
      <TableContainer component={Paper} elevation={0} sx={{ mb: 2, border: '1px solid #e0e0e0', borderRadius: 1 }}>
        <Table size="small">{children}</Table>
      </TableContainer>
    ),
    thead: ({ children }: any) => <TableHead sx={{ backgroundColor: accentColor }}>{children}</TableHead>,
    tbody: ({ children }: any) => <TableBody>{children}</TableBody>,
    tr: ({ children }: any) => <TableRow>{children}</TableRow>,
    th: ({ children }: any) => <TableCell sx={{ color: '#fff', fontWeight: 600, fontSize: 13, borderColor: 'rgba(255,255,255,0.2)' }}>{children}</TableCell>,
    td: ({ children }: any) => <TableCell sx={{ fontSize: 13 }}>{children}</TableCell>,
    li: ({ children }: any) => <li style={{ marginBottom: 4 }}>{children}</li>,
  };

  // ── Render ──────────────────────────────────────────────────────────────────

  return (
    <div style={{ display: 'flex', height: '100vh', fontFamily: 'system-ui, sans-serif', background: '#f5f5f5' }}>
      {/* Sidebar */}
      <div style={{
        width: sidebarOpen ? 280 : 0,
        minWidth: sidebarOpen ? 280 : 0,
        overflow: 'hidden',
        transition: 'width 0.2s, min-width 0.2s',
        background: '#1a1a2e',
        color: '#fff',
        display: 'flex',
        flexDirection: 'column',
        flexShrink: 0,
      }}>
        <div style={{ padding: '20px 16px 12px', borderBottom: '1px solid #333' }}>
          <div style={{ fontSize: 11, letterSpacing: 1, color: '#aaa', textTransform: 'uppercase', marginBottom: 4 }}>AQA A-Level</div>
          <div style={{ fontSize: 18, fontWeight: 700 }}>Physics Revision</div>
        </div>
        <div style={{ overflowY: 'auto', flex: 1, padding: '8px 0' }}>
          {APHY_MODULES.map(mod => (
            <div key={mod.id}>
              <div style={{
                padding: '10px 16px 6px',
                fontSize: 11,
                fontWeight: 700,
                letterSpacing: 0.5,
                textTransform: 'uppercase',
                color: MODULE_COLORS[mod.id] || '#aaa',
                marginTop: 4,
              }}>
                {mod.name}
              </div>
              {mod.lessons.map(lesson => (
                <button
                  key={lesson.id}
                  onClick={() => handleLessonSelect(mod, lesson)}
                  style={{
                    display: 'block',
                    width: '100%',
                    textAlign: 'left',
                    background: selectedLesson.id === lesson.id ? 'rgba(255,255,255,0.12)' : 'transparent',
                    border: 'none',
                    borderLeft: selectedLesson.id === lesson.id ? `3px solid ${MODULE_COLORS[mod.id]}` : '3px solid transparent',
                    color: selectedLesson.id === lesson.id ? '#fff' : '#bbb',
                    padding: '7px 16px 7px 14px',
                    fontSize: 13,
                    cursor: 'pointer',
                    lineHeight: 1.3,
                  }}
                >
                  {lesson.name}
                </button>
              ))}
            </div>
          ))}
        </div>
      </div>

      {/* Main content */}
      <div style={{ flex: 1, display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
        {/* Top bar */}
        <div style={{
          background: accentColor,
          color: '#fff',
          padding: '10px 20px',
          display: 'flex',
          alignItems: 'center',
          gap: 12,
          flexShrink: 0,
          flexWrap: 'wrap',
        }}>
          <button
            onClick={() => setSidebarOpen(o => !o)}
            style={{ background: 'rgba(255,255,255,0.2)', border: 'none', color: '#fff', borderRadius: 4, padding: '4px 10px', cursor: 'pointer', fontSize: 16 }}
          >
            ☰
          </button>
          <div style={{ flex: 1 }}>
            <div style={{ fontSize: 11, opacity: 0.8 }}>{selectedModule.name}</div>
            <div style={{ fontSize: 16, fontWeight: 600 }}>{selectedLesson.name}</div>
          </div>
          {/* Mode buttons */}
          <div style={{ display: 'flex', gap: 6 }}>
            {(['read', 'speedread', 'audio', 'flashcards'] as Mode[]).map(m => (
              <button
                key={m}
                onClick={() => setMode(m)}
                style={{
                  background: mode === m ? 'rgba(255,255,255,0.95)' : 'rgba(255,255,255,0.2)',
                  color: mode === m ? accentColor : '#fff',
                  border: 'none',
                  borderRadius: 16,
                  padding: '5px 14px',
                  fontSize: 12,
                  fontWeight: 600,
                  cursor: 'pointer',
                  textTransform: 'capitalize',
                }}
              >
                {m === 'speedread' ? '⚡ Speed' : m === 'audio' ? '🔊 Audio' : m === 'flashcards' ? '🃏 Cards' : '📖 Read'}
              </button>
            ))}
          </div>
        </div>

        {/* ── READ mode ────────────────────────────────────────────────────────── */}
        {mode === 'read' && (
          <div style={{ flex: 1, overflowY: 'auto', padding: '24px 32px', maxWidth: 800, margin: '0 auto', width: '100%', boxSizing: 'border-box' }}>
            <div style={{
              background: '#fff',
              borderRadius: 8,
              padding: '28px 32px',
              boxShadow: '0 1px 4px rgba(0,0,0,0.08)',
              lineHeight: 1.7,
              fontSize: 15,
            }}>
              <ReactMarkdown
                remarkPlugins={[remarkMath, remarkGfm]}
                rehypePlugins={[rehypeKatex]}
                components={mdComponents}
              >
                {selectedLesson.content}
              </ReactMarkdown>
            </div>
            <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: 20, gap: 12 }}>
              {prevLesson ? (
                <button onClick={() => handleLessonSelect(selectedModule, prevLesson)} style={navBtnStyle(accentColor)}>
                  ← {prevLesson.name}
                </button>
              ) : <div />}
              {nextLesson ? (
                <button onClick={() => handleLessonSelect(selectedModule, nextLesson)} style={navBtnStyle(accentColor)}>
                  {nextLesson.name} →
                </button>
              ) : <div />}
            </div>
          </div>
        )}

        {/* ── SPEED READ mode ──────────────────────────────────────────────────── */}
        {mode === 'speedread' && (
          <div style={{ flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: 24, padding: 24 }}>
            {/* WPM selector */}
            <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
              {[150, 300, 450].map(w => (
                <button key={w} onClick={() => setSrWpm(w)} style={{
                  background: srWpm === w ? accentColor : '#e0e0e0',
                  color: srWpm === w ? '#fff' : '#333',
                  border: 'none', borderRadius: 20, padding: '4px 14px', cursor: 'pointer', fontWeight: 600,
                }}>
                  {w} WPM
                </button>
              ))}
            </div>

            {/* Word display */}
            <div style={{
              background: '#1a1a2e',
              borderRadius: 12,
              width: '100%',
              maxWidth: 600,
              minHeight: 160,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              flexDirection: 'column',
              gap: 8,
            }}>
              {srComplete ? (
                <div style={{ color: '#aaa', fontSize: 18 }}>✓ Complete</div>
              ) : (
                <>
                  <div style={{ color: '#fff', fontSize: 36, fontWeight: 700, letterSpacing: 1 }}>
                    {srWords[srIndex] || '—'}
                  </div>
                  <div style={{ color: '#555', fontSize: 13 }}>{srIndex + 1} / {srWords.length}</div>
                </>
              )}
            </div>

            {/* Progress bar */}
            <div style={{ width: '100%', maxWidth: 600, height: 6, background: '#e0e0e0', borderRadius: 3 }}>
              <div style={{ width: `${srWords.length ? (srIndex / srWords.length) * 100 : 0}%`, height: '100%', background: accentColor, borderRadius: 3, transition: 'width 0.1s' }} />
            </div>

            {/* Controls */}
            <div style={{ display: 'flex', gap: 12 }}>
              {!srPlaying && !srComplete && (
                <button onClick={() => { setSrIndex(0); setSrPlaying(true); }} style={srBtnStyle(accentColor)}>
                  {srIndex === 0 ? '▶ Start' : '▶ Resume'}
                </button>
              )}
              {srPlaying && (
                <button onClick={() => setSrPlaying(false)} style={srBtnStyle('#555')}>⏸ Pause</button>
              )}
              {(srIndex > 0 || srComplete) && (
                <button onClick={() => { setSrIndex(0); setSrPlaying(false); setSrComplete(false); }} style={srBtnStyle('#888')}>↺ Restart</button>
              )}
            </div>
          </div>
        )}

        {/* ── AUDIO mode ───────────────────────────────────────────────────────── */}
        {mode === 'audio' && (
          <div style={{ flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: 24, padding: 24 }}>
            <div style={{ fontSize: 48 }}>{audioPlaying ? '🔊' : '🔇'}</div>
            <div style={{ fontSize: 18, fontWeight: 600, color: '#333', textAlign: 'center', maxWidth: 400 }}>
              {audioPlaying ? 'Reading aloud…' : 'Press play to listen to this lesson'}
            </div>
            <div style={{ display: 'flex', gap: 12 }}>
              {!audioPlaying ? (
                <button onClick={handleAudioPlay} style={srBtnStyle(accentColor)}>▶ Play</button>
              ) : (
                <button onClick={handleAudioStop} style={srBtnStyle('#c62828')}>⏹ Stop</button>
              )}
            </div>
            <div style={{ color: '#888', fontSize: 12, maxWidth: 400, textAlign: 'center' }}>
              Uses your device's text-to-speech. Equations are read as "equation". Works best with earphones.
            </div>
          </div>
        )}

        {/* ── FLASHCARDS mode ──────────────────────────────────────────────────── */}
        {mode === 'flashcards' && (
          <div style={{ flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: 20, padding: 24 }}>
            {flashcards.length === 0 ? (
              <div style={{ color: '#888', fontSize: 16 }}>No flashcards found for this lesson.</div>
            ) : (
              <>
                <div style={{ color: '#888', fontSize: 13 }}>Card {fcIndex + 1} of {flashcards.length} — tap to flip</div>

                {/* Card */}
                <div
                  onClick={() => setFcFlipped(f => !f)}
                  style={{
                    width: '100%',
                    maxWidth: 560,
                    minHeight: fcFlipped ? undefined : 220,
                    background: fcFlipped ? '#f0f4ff' : '#fff',
                    color: '#1a1a2e',
                    borderRadius: 14,
                    boxShadow: fcFlipped ? `0 4px 20px ${accentColor}40` : '0 4px 20px rgba(0,0,0,0.12)',
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: fcFlipped ? 'stretch' : 'center',
                    justifyContent: fcFlipped ? 'flex-start' : 'center',
                    padding: 28,
                    cursor: 'pointer',
                    transition: 'background 0.25s, box-shadow 0.25s',
                    textAlign: fcFlipped ? 'left' : 'center',
                  }}
                >
                  <div style={{ fontSize: 11, fontWeight: 700, letterSpacing: 1, color: accentColor, textTransform: 'uppercase', marginBottom: 12 }}>
                    {fcFlipped ? 'Answer' : 'Question'}
                  </div>
                  {fcFlipped ? (
                    <div style={{ fontSize: 14, lineHeight: 1.6 }}>
                      <ReactMarkdown
                        remarkPlugins={[remarkMath, remarkGfm]}
                        rehypePlugins={[rehypeKatex]}
                        components={mdComponents}
                      >
                        {fcCard.back}
                      </ReactMarkdown>
                    </div>
                  ) : (
                    <div style={{ fontSize: 20, fontWeight: 700, lineHeight: 1.5 }}>
                      {fcCard.front}
                    </div>
                  )}
                </div>

                {/* Nav */}
                <div style={{ display: 'flex', gap: 12 }}>
                  <button
                    disabled={fcIndex === 0}
                    onClick={() => { setFcIndex(i => i - 1); setFcFlipped(false); }}
                    style={srBtnStyle(fcIndex === 0 ? '#ccc' : '#555')}
                  >
                    ← Prev
                  </button>
                  <button
                    onClick={() => { setFcFlipped(false); setFcIndex(0); }}
                    style={srBtnStyle('#888')}
                  >
                    ↺ Restart
                  </button>
                  <button
                    disabled={fcIndex === flashcards.length - 1}
                    onClick={() => { setFcIndex(i => i + 1); setFcFlipped(false); }}
                    style={srBtnStyle(fcIndex === flashcards.length - 1 ? '#ccc' : accentColor)}
                  >
                    Next →
                  </button>
                </div>
              </>
            )}
          </div>
        )}
      </div>
    </div>
  );
}

function navBtnStyle(color: string): React.CSSProperties {
  return {
    background: color,
    color: '#fff',
    border: 'none',
    borderRadius: 6,
    padding: '10px 18px',
    cursor: 'pointer',
    fontSize: 13,
    maxWidth: '45%',
    textAlign: 'center',
    lineHeight: 1.3,
  };
}

function srBtnStyle(color: string): React.CSSProperties {
  return {
    background: color,
    color: '#fff',
    border: 'none',
    borderRadius: 8,
    padding: '10px 22px',
    cursor: 'pointer',
    fontSize: 14,
    fontWeight: 600,
  };
}
