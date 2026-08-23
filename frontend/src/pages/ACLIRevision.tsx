import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkMath from 'remark-math';
import remarkGfm from 'remark-gfm';
import rehypeKatex from 'rehype-katex';
import 'katex/dist/katex.min.css';
import { ACLI_GLOSSARY, ACLI_MODULES, ACLILesson, ACLIModule } from './acliContent';

const MODULE_COLORS: Record<number, string> = {
  5001: '#0f4c81',
  5002: '#17643a',
  5003: '#7a3e00',
};

type Mode = 'read' | 'speedread' | 'audio' | 'flashcards';

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
        .replace(/^\s*[-*]\s+/gm, '• ')
        .trim();
      if (back.length > 20) cards.push({ front: currentHeading, back });
    }
  };

  for (const line of lines) {
    const h2 = line.match(/^## (.+)/);
    const h1 = line.match(/^# (.+)/);
    if (h1 || h2) {
      flush();
      currentHeading = (h1 || h2)![1];
      bodyLines = [];
    } else {
      bodyLines.push(line);
    }
  }
  flush();
  return cards;
}

function materializeGlossary(content: string, includeExplainers: boolean): string {
  return content.replace(/\[\[([^|\]]+)\|([^\]]+)\]\]/g, (_, term: string, micro: string) => {
    if (!includeExplainers) return term;
    return `${term} (${micro})`;
  });
}

function stripMnemonicBlocks(md: string): string {
  return md.replace(/^##\s+Mnemonic[\s\S]*?(?=^##\s+|^#\s+|$)/gim, '').trim();
}

function markdownToPlainText(md: string): string {
  return md
    .replace(/```[\s\S]*?```/g, '')
    .replace(/\$\$[\s\S]*?\$\$/g, 'equation')
    .replace(/\$[^$]+\$/g, 'equation')
    .replace(/^\|.*\|$/gm, '')
    .replace(/^#+\s+/gm, '')
    .replace(/\*\*([^*]+)\*\*/g, '$1')
    .replace(/\*([^*]+)\*/g, '$1')
    .replace(/`([^`]+)`/g, '$1')
    .replace(/^\s*[-*]\s+/gm, '')
    .replace(/\n{3,}/g, '\n\n')
    .trim();
}

const ACLIRevision: React.FC = () => {
  const [selectedModule, setSelectedModule] = useState<ACLIModule>(ACLI_MODULES[0]);
  const [selectedLesson, setSelectedLesson] = useState<ACLILesson>(ACLI_MODULES[0].lessons[0]);
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [mode, setMode] = useState<Mode>('read');
  const [includeExplainers, setIncludeExplainers] = useState(false);
  const [examMode, setExamMode] = useState(false);
  const [revealHints, setRevealHints] = useState(false);

  const [srWords, setSrWords] = useState<string[]>([]);
  const [srIndex, setSrIndex] = useState(0);
  const [srPlaying, setSrPlaying] = useState(false);
  const [srComplete, setSrComplete] = useState(false);
  const [srWpm, setSrWpm] = useState(280);
  const srTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const [audioPlaying, setAudioPlaying] = useState(false);

  const [flashcards, setFlashcards] = useState<Flashcard[]>([]);
  const [fcIndex, setFcIndex] = useState(0);
  const [fcFlipped, setFcFlipped] = useState(false);

  const accentColor = MODULE_COLORS[selectedModule.id] || '#0f4c81';

  const renderedLesson = useMemo(() => {
    const hintsVisible = !examMode || revealHints;
    const glossaryEnabled = includeExplainers && hintsVisible;
    const mnemonicVisible = hintsVisible;

    let content = materializeGlossary(selectedLesson.content, glossaryEnabled);
    if (!mnemonicVisible) content = stripMnemonicBlocks(content);
    return content;
  }, [selectedLesson.content, includeExplainers, examMode, revealHints]);

  const handleLessonSelect = useCallback((mod: ACLIModule, lesson: ACLILesson) => {
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
    setRevealHints(false);
    if (window.innerWidth < 768) setSidebarOpen(false);
  }, [audioPlaying]);

  useEffect(() => {
    if (!examMode) setRevealHints(false);
  }, [examMode]);

  useEffect(() => {
    if (mode === 'speedread') {
      const text = markdownToPlainText(renderedLesson);
      const words = text.split(/\s+/).filter(Boolean);
      setSrWords(words);
      setSrIndex(0);
      setSrPlaying(false);
      setSrComplete(false);
    }
  }, [mode, renderedLesson]);

  useEffect(() => {
    if (!srPlaying || srIndex >= srWords.length) return;
    const ms = 60000 / srWpm;
    srTimerRef.current = setTimeout(() => {
      if (srIndex < srWords.length - 1) {
        setSrIndex((i) => i + 1);
      } else {
        setSrPlaying(false);
        setSrComplete(true);
      }
    }, ms);

    return () => {
      if (srTimerRef.current) clearTimeout(srTimerRef.current);
    };
  }, [srPlaying, srIndex, srWpm, srWords.length]);

  useEffect(() => {
    if (mode === 'flashcards') {
      setFlashcards(extractFlashcards(renderedLesson));
      setFcIndex(0);
      setFcFlipped(false);
    }
  }, [mode, renderedLesson]);

  const currentIdx = selectedModule.lessons.findIndex((l) => l.id === selectedLesson.id);
  const prevLesson = currentIdx > 0 ? selectedModule.lessons[currentIdx - 1] : null;
  const nextLesson = currentIdx < selectedModule.lessons.length - 1 ? selectedModule.lessons[currentIdx + 1] : null;

  const handleAudioPlay = () => {
    const text = markdownToPlainText(renderedLesson);
    const utt = new SpeechSynthesisUtterance(text);
    utt.rate = 0.95;
    utt.onend = () => setAudioPlaying(false);
    window.speechSynthesis.speak(utt);
    setAudioPlaying(true);
  };

  const handleAudioStop = () => {
    window.speechSynthesis?.cancel();
    setAudioPlaying(false);
  };

  const glossaryEntries = Object.entries(ACLI_GLOSSARY);
  const fcCard = flashcards[fcIndex];

  const mdComponents: any = {
    h1: ({ children }: any) => <h1 style={{ fontSize: 24, fontWeight: 700, color: accentColor, marginBottom: 14, marginTop: 0 }}>{children}</h1>,
    h2: ({ children }: any) => <h2 style={{ fontSize: 17, fontWeight: 600, marginTop: 22, marginBottom: 8, borderBottom: `2px solid ${accentColor}26`, paddingBottom: 4 }}>{children}</h2>,
    h3: ({ children }: any) => <h3 style={{ fontSize: 15, fontWeight: 600, marginTop: 14, marginBottom: 6 }}>{children}</h3>,
    strong: ({ children }: any) => <strong style={{ color: accentColor }}>{children}</strong>,
    code: ({ inline, children }: any) => inline ? (
      <code style={{ background: '#f2f4f7', padding: '1px 4px', borderRadius: 3 }}>{children}</code>
    ) : (
      <pre style={{ background: '#1f2937', color: '#e5e7eb', padding: '12px 14px', borderRadius: 6, overflowX: 'auto' }}><code>{children}</code></pre>
    ),
    li: ({ children }: any) => <li style={{ marginBottom: 4 }}>{children}</li>,
  };

  const modeButton = (active: boolean, color: string): React.CSSProperties => ({
    border: active ? `2px solid ${color}` : '1px solid #d0d7de',
    background: active ? `${color}1A` : '#fff',
    color: '#111827',
    borderRadius: 8,
    padding: '8px 10px',
    cursor: 'pointer',
    fontWeight: 600,
    fontSize: 13,
  });

  return (
    <div style={{ display: 'flex', height: '100vh', background: '#f3f5f7', fontFamily: 'system-ui, sans-serif' }}>
      <div
        style={{
          width: sidebarOpen ? 330 : 0,
          minWidth: sidebarOpen ? 330 : 0,
          overflow: 'hidden',
          transition: 'width 0.2s, min-width 0.2s',
          background: '#ffffff',
          borderRight: '1px solid #e5e7eb',
        }}
      >
        <div style={{ padding: 16, borderBottom: '1px solid #e5e7eb' }}>
          <div style={{ fontSize: 15, fontWeight: 700, color: accentColor }}>ACLI Clinical Immunology Modules</div>
          <div style={{ fontSize: 12, color: '#6b7280', marginTop: 4 }}>FRCPath-focused adaptive revision content</div>
        </div>

        <div style={{ padding: 12, overflowY: 'auto', maxHeight: 'calc(100vh - 74px)' }}>
          {ACLI_MODULES.map((mod) => (
            <div key={mod.id} style={{ marginBottom: 16 }}>
              <div style={{ fontSize: 12, fontWeight: 700, color: MODULE_COLORS[mod.id] || '#6b7280', margin: '0 8px 8px' }}>
                {mod.name}
              </div>
              {mod.lessons.map((lesson) => {
                const selected = selectedLesson.id === lesson.id;
                return (
                  <button
                    key={lesson.id}
                    onClick={() => handleLessonSelect(mod, lesson)}
                    style={{
                      width: '100%',
                      textAlign: 'left',
                      border: 'none',
                      borderLeft: selected ? `4px solid ${MODULE_COLORS[mod.id] || accentColor}` : '4px solid transparent',
                      borderRadius: 6,
                      marginBottom: 6,
                      padding: '10px 10px',
                      background: selected ? '#eef2ff' : '#fff',
                      color: '#111827',
                      cursor: 'pointer',
                    }}
                  >
                    <div style={{ fontSize: 13, fontWeight: selected ? 700 : 500 }}>{lesson.title}</div>
                  </button>
                );
              })}
            </div>
          ))}
        </div>
      </div>

      <div style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
        <div style={{ background: accentColor, color: '#fff', padding: '12px 16px', display: 'flex', alignItems: 'center', gap: 10 }}>
          <button onClick={() => setSidebarOpen((s) => !s)} style={{ border: 'none', borderRadius: 6, background: '#ffffff26', color: '#fff', padding: '6px 10px', cursor: 'pointer' }}>
            {sidebarOpen ? 'Hide' : 'Show'} module list
          </button>
          <div style={{ fontSize: 17, fontWeight: 700, flex: 1 }}>{selectedLesson.title}</div>
          <label style={{ display: 'flex', alignItems: 'center', gap: 6, fontSize: 12 }}>
            <input
              type="checkbox"
              checked={examMode}
              onChange={(e) => setExamMode(e.target.checked)}
            />
            Exam mode
          </label>
          {examMode && (
            <button
              onClick={() => setRevealHints((v) => !v)}
              style={{ border: 'none', borderRadius: 6, background: '#ffffff26', color: '#fff', padding: '6px 10px', cursor: 'pointer', fontSize: 12 }}
            >
              {revealHints ? 'Hide hints' : 'Reveal hints'}
            </button>
          )}
          <label style={{ display: 'flex', alignItems: 'center', gap: 6, fontSize: 12 }}>
            <input
              type="checkbox"
              checked={includeExplainers}
              onChange={(e) => setIncludeExplainers(e.target.checked)}
              disabled={examMode && !revealHints}
            />
            Micro explainers
          </label>
        </div>

        <div style={{ padding: 12, background: '#fff', borderBottom: '1px solid #e5e7eb', display: 'flex', alignItems: 'center', gap: 8, flexWrap: 'wrap' }}>
          <button style={modeButton(mode === 'read', accentColor)} onClick={() => setMode('read')}>Read</button>
          <button style={modeButton(mode === 'speedread', accentColor)} onClick={() => setMode('speedread')}>Speedread</button>
          <button style={modeButton(mode === 'audio', accentColor)} onClick={() => setMode('audio')}>Audio</button>
          <button style={modeButton(mode === 'flashcards', accentColor)} onClick={() => setMode('flashcards')}>Flashcards</button>
          {mode === 'speedread' && (
            <>
              <span style={{ fontSize: 12, color: '#4b5563', marginLeft: 8 }}>WPM</span>
              <input
                type="range"
                min={180}
                max={520}
                step={20}
                value={srWpm}
                onChange={(e) => setSrWpm(Number(e.target.value))}
              />
              <span style={{ fontSize: 12, color: '#111827', minWidth: 36 }}>{srWpm}</span>
            </>
          )}
        </div>

        <div style={{ flex: 1, overflowY: 'auto', padding: 20 }}>
          {mode === 'read' && (
            <div style={{ display: 'grid', gridTemplateColumns: 'minmax(0, 1fr) 280px', gap: 16 }}>
              <div style={{ background: '#fff', borderRadius: 10, padding: 20, boxShadow: '0 1px 2px rgba(0,0,0,0.08)' }}>
                <ReactMarkdown
                  remarkPlugins={[remarkMath, remarkGfm]}
                  rehypePlugins={[rehypeKatex]}
                  components={mdComponents}
                >
                  {renderedLesson}
                </ReactMarkdown>
              </div>
              <div style={{ background: '#fff', borderRadius: 10, padding: 14, boxShadow: '0 1px 2px rgba(0,0,0,0.08)' }}>
                <div style={{ fontSize: 13, fontWeight: 700, marginBottom: 8, color: accentColor }}>Quick glossary</div>
                <div style={{ fontSize: 12, color: '#6b7280', marginBottom: 10 }}>
                  {examMode && !revealHints
                    ? 'Exam mode is active: glossary hints are hidden until Reveal hints is toggled.'
                    : 'Toggle "Micro explainers" to inject these directly into lesson text.'}
                </div>
                <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                  {glossaryEntries.slice(0, 10).map(([term, definition]) => (
                    <div key={term} style={{ border: '1px solid #e5e7eb', borderRadius: 8, padding: 8 }}>
                      <div style={{ fontSize: 12, fontWeight: 700 }}>{term}</div>
                      <div style={{ fontSize: 12, color: '#4b5563' }}>{definition}</div>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          )}

          {mode === 'speedread' && (
            <div style={{ maxWidth: 860, margin: '0 auto', textAlign: 'center', paddingTop: 24 }}>
              <div style={{ color: '#6b7280', fontSize: 13, marginBottom: 8 }}>
                Word {Math.min(srIndex + 1, Math.max(srWords.length, 1))} of {srWords.length}
              </div>
              <div style={{ fontSize: 52, fontWeight: 700, minHeight: 82, color: accentColor }}>
                {srWords[srIndex] || 'Ready'}
              </div>
              <div style={{ marginTop: 24, display: 'flex', justifyContent: 'center', gap: 10 }}>
                {!srPlaying ? (
                  <button onClick={() => setSrPlaying(true)} style={modeButton(true, accentColor)}>
                    {srComplete ? 'Restart' : 'Play'}
                  </button>
                ) : (
                  <button onClick={() => setSrPlaying(false)} style={modeButton(true, '#b91c1c')}>Pause</button>
                )}
                <button
                  onClick={() => {
                    setSrPlaying(false);
                    setSrIndex(0);
                    setSrComplete(false);
                  }}
                  style={modeButton(false, accentColor)}
                >
                  Reset
                </button>
              </div>
            </div>
          )}

          {mode === 'audio' && (
            <div style={{ maxWidth: 760, margin: '0 auto', textAlign: 'center', paddingTop: 24 }}>
              <div style={{ fontSize: 14, color: '#6b7280', marginBottom: 8 }}>
                Device text-to-speech reader for this lesson.
              </div>
              <div style={{ fontSize: 14, color: '#4b5563', marginBottom: 16 }}>
                Micro explainers are currently {(includeExplainers && (!examMode || revealHints)) ? 'ON' : 'OFF'}.
              </div>
              <div style={{ display: 'flex', justifyContent: 'center', gap: 10 }}>
                {!audioPlaying ? (
                  <button onClick={handleAudioPlay} style={modeButton(true, accentColor)}>Play</button>
                ) : (
                  <button onClick={handleAudioStop} style={modeButton(true, '#b91c1c')}>Stop</button>
                )}
              </div>
              <div style={{ marginTop: 20, background: '#fff', borderRadius: 10, padding: 14, textAlign: 'left', boxShadow: '0 1px 2px rgba(0,0,0,0.08)' }}>
                <div style={{ fontSize: 13, fontWeight: 700, marginBottom: 8 }}>Audio script preview</div>
                <div style={{ fontSize: 13, color: '#374151', whiteSpace: 'pre-wrap' }}>
                  {markdownToPlainText(renderedLesson).slice(0, 1200)}
                  {markdownToPlainText(renderedLesson).length > 1200 ? ' ...' : ''}
                </div>
              </div>
            </div>
          )}

          {mode === 'flashcards' && (
            <div style={{ maxWidth: 860, margin: '0 auto', textAlign: 'center', paddingTop: 24 }}>
              {flashcards.length === 0 ? (
                <div style={{ color: '#6b7280' }}>No flashcards found for this lesson.</div>
              ) : (
                <>
                  <div style={{ fontSize: 13, color: '#6b7280', marginBottom: 10 }}>Card {fcIndex + 1} of {flashcards.length}</div>
                  <button
                    onClick={() => setFcFlipped((f) => !f)}
                    style={{
                      width: '100%',
                      textAlign: 'left',
                      borderRadius: 12,
                      border: `1px solid ${accentColor}66`,
                      background: '#fff',
                      padding: 16,
                      cursor: 'pointer',
                      boxShadow: '0 1px 2px rgba(0,0,0,0.08)',
                    }}
                  >
                    <div style={{ fontSize: 12, color: '#6b7280', marginBottom: 8 }}>{fcFlipped ? 'Back' : 'Front'} (click to flip)</div>
                    <div style={{ fontSize: 18, fontWeight: 700, color: accentColor, marginBottom: 12 }}>{fcCard.front}</div>
                    {fcFlipped && <div style={{ fontSize: 14, color: '#111827', whiteSpace: 'pre-wrap' }}>{fcCard.back}</div>}
                  </button>

                  <div style={{ marginTop: 12, display: 'flex', justifyContent: 'space-between' }}>
                    <button
                      onClick={() => {
                        setFcIndex((i) => Math.max(0, i - 1));
                        setFcFlipped(false);
                      }}
                      disabled={fcIndex === 0}
                      style={modeButton(fcIndex > 0, accentColor)}
                    >
                      Previous
                    </button>
                    <button
                      onClick={() => {
                        setFcIndex((i) => Math.min(flashcards.length - 1, i + 1));
                        setFcFlipped(false);
                      }}
                      disabled={fcIndex === flashcards.length - 1}
                      style={modeButton(fcIndex < flashcards.length - 1, accentColor)}
                    >
                      Next
                    </button>
                  </div>
                </>
              )}
            </div>
          )}
        </div>

        <div style={{ padding: 10, borderTop: '1px solid #e5e7eb', background: '#fff', display: 'flex', justifyContent: 'space-between' }}>
          <button
            disabled={!prevLesson}
            onClick={() => prevLesson && handleLessonSelect(selectedModule, prevLesson)}
            style={modeButton(Boolean(prevLesson), accentColor)}
          >
            Previous lesson
          </button>
          <button
            disabled={!nextLesson}
            onClick={() => nextLesson && handleLessonSelect(selectedModule, nextLesson)}
            style={modeButton(Boolean(nextLesson), accentColor)}
          >
            Next lesson
          </button>
        </div>
      </div>
    </div>
  );
};

export default ACLIRevision;
