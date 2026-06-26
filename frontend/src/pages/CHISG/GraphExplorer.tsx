import React, { useCallback, useEffect, useRef, useState } from 'react';
import ForceGraph2D, { ForceGraphMethods } from 'react-force-graph-2d';
import {
  Box,
  Button,
  Chip,
  CircularProgress,
  Divider,
  FormControl,
  IconButton,
  InputLabel,
  MenuItem,
  Paper,
  Select,
  Stack,
  Tooltip,
  Typography,
  alpha,
  useTheme,
} from '@mui/material';
import {
  Close as CloseIcon,
  FilterList as FilterIcon,
  Hub as GraphIcon,
  Refresh as RefreshIcon,
} from '@mui/icons-material';
import { Link } from 'react-router-dom';
import { apiClient, API_ENDPOINTS } from '../../config/api';
import SEO from '../../components/SEO';

// ── Types ────────────────────────────────────────────────────────────────────

interface GraphNode {
  id: string;
  label: string;
  domain: string;
  val: number;
  // runtime fields added by force-graph
  x?: number;
  y?: number;
  fx?: number;
  fy?: number;
}

interface GraphEdge {
  source: string | GraphNode;
  target: string | GraphNode;
  relation: string;
  inverse_relation?: string;
  domain?: string;
  source_title?: string;
  evidence_context?: string;
  confidence: number;
}

interface GraphData {
  nodes: GraphNode[];
  edges: GraphEdge[];
  meta: { node_count: number; edge_count: number };
}

// ── Colour palette ───────────────────────────────────────────────────────────

const DOMAIN_COLOURS: Record<string, string> = {
  microbiology: '#4ade80',
  immunology: '#60a5fa',
  education: '#f472b6',
  biology: '#34d399',
  physics: '#a78bfa',
  chemistry: '#fb923c',
  psychology: '#f87171',
  manual: '#94a3b8',
  default: '#cbd5e1',
};

const RELATION_COLOURS: Record<string, string> = {
  causes: '#ef4444',
  inhibits: '#f97316',
  enables: '#22c55e',
  'is composed of': '#06b6d4',
  'is a type of': '#8b5cf6',
  'has property': '#d946ef',
  'is used for': '#0ea5e9',
  determines: '#eab308',
  represents: '#64748b',
  'correlates with': '#14b8a6',
  supports: '#84cc16',
  refutes: '#dc2626',
  regulates: '#f59e0b',
  'is a mechanism of': '#7c3aed',
  'is a marker for': '#be185d',
  'interacts with': '#0891b2',
  'is found in': '#15803d',
  'is an example of': '#6366f1',
  contradicts: '#b91c1c',
  default: '#475569',
};

function domainColour(domain: string): string {
  return DOMAIN_COLOURS[domain?.toLowerCase()] ?? DOMAIN_COLOURS.default;
}

function relationColour(relation: string): string {
  const key = relation?.toLowerCase().trim();
  return RELATION_COLOURS[key] ?? RELATION_COLOURS.default;
}

// ── Component ────────────────────────────────────────────────────────────────

const GraphExplorer: React.FC = () => {
  const theme = useTheme();
  const fgRef = useRef<ForceGraphMethods>();
  const containerRef = useRef<HTMLDivElement>(null);
  const [dimensions, setDimensions] = useState({ width: 800, height: 600 });

  const [graphData, setGraphData] = useState<GraphData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Filters
  const [domainFilter, setDomainFilter] = useState('');
  const [relationFilter, setRelationFilter] = useState('');

  // Selection
  const [selectedNode, setSelectedNode] = useState<GraphNode | null>(null);
  const [selectedEdge, setSelectedEdge] = useState<GraphEdge | null>(null);
  const [highlightNodes, setHighlightNodes] = useState<Set<string>>(new Set());
  const [highlightEdges, setHighlightEdges] = useState<Set<GraphEdge>>(new Set());

  // ── Fetch ──────────────────────────────────────────────────────────────────

  const fetchGraph = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const params = new URLSearchParams({ limit: '800' });
      if (domainFilter) params.set('domain', domainFilter);
      const data = await apiClient.get<GraphData>(`${API_ENDPOINTS.chisg.graph}?${params}`);
      setGraphData(data);
    } catch (e: any) {
      setError(e?.message ?? 'Failed to load graph data');
    } finally {
      setLoading(false);
    }
  }, [domainFilter]);

  useEffect(() => { fetchGraph(); }, [fetchGraph]);

  // ── Responsive dimensions ──────────────────────────────────────────────────

  useEffect(() => {
    const update = () => {
      if (containerRef.current) {
        setDimensions({
          width: containerRef.current.clientWidth,
          height: containerRef.current.clientHeight,
        });
      }
    };
    update();
    const ro = new ResizeObserver(update);
    if (containerRef.current) ro.observe(containerRef.current);
    return () => ro.disconnect();
  }, []);

  // ── Compute filtered view ──────────────────────────────────────────────────

  const filteredGraphData = React.useMemo(() => {
    if (!graphData) return { nodes: [], links: [] };
    let edges = graphData.edges;
    if (relationFilter) {
      edges = edges.filter(e => e.relation?.toLowerCase() === relationFilter.toLowerCase());
    }
    const usedIds = new Set<string>();
    edges.forEach(e => {
      const src = typeof e.source === 'string' ? e.source : (e.source as GraphNode).id;
      const tgt = typeof e.target === 'string' ? e.target : (e.target as GraphNode).id;
      usedIds.add(src);
      usedIds.add(tgt);
    });
    const nodes = graphData.nodes.filter(n => usedIds.has(n.id));
    return { nodes, links: edges };
  }, [graphData, relationFilter]);

  // ── Highlight helpers ──────────────────────────────────────────────────────

  const clearHighlight = useCallback(() => {
    setHighlightNodes(new Set());
    setHighlightEdges(new Set());
    setSelectedNode(null);
    setSelectedEdge(null);
  }, []);

  const handleNodeClick = useCallback((node: GraphNode) => {
    setSelectedNode(node);
    setSelectedEdge(null);
    const connected = new Set<string>([node.id]);
    const connEdges = new Set<GraphEdge>();
    filteredGraphData.links.forEach((e: GraphEdge) => {
      const src = typeof e.source === 'string' ? e.source : (e.source as GraphNode).id;
      const tgt = typeof e.target === 'string' ? e.target : (e.target as GraphNode).id;
      if (src === node.id || tgt === node.id) {
        connected.add(src);
        connected.add(tgt);
        connEdges.add(e);
      }
    });
    setHighlightNodes(connected);
    setHighlightEdges(connEdges);
    if (fgRef.current) {
      fgRef.current.centerAt(node.x, node.y, 600);
      fgRef.current.zoom(3, 600);
    }
  }, [filteredGraphData]);

  const handleEdgeClick = useCallback((edge: GraphEdge) => {
    setSelectedEdge(edge);
    setSelectedNode(null);
    const src = typeof edge.source === 'string' ? edge.source : (edge.source as GraphNode).id;
    const tgt = typeof edge.target === 'string' ? edge.target : (edge.target as GraphNode).id;
    setHighlightNodes(new Set([src, tgt]));
    setHighlightEdges(new Set([edge]));
  }, []);

  // ── Derived options ────────────────────────────────────────────────────────

  const availableDomains = React.useMemo(() => {
    if (!graphData) return [];
    return Array.from(new Set(graphData.nodes.map(n => n.domain).filter(Boolean))).sort();
  }, [graphData]);

  const availableRelations = React.useMemo(() => {
    if (!graphData) return [];
    return Array.from(new Set(graphData.edges.map(e => e.relation).filter(Boolean))).sort();
  }, [graphData]);

  // ── Render helpers ─────────────────────────────────────────────────────────

  const nodeCanvasObject = useCallback(
    (node: GraphNode, ctx: CanvasRenderingContext2D, globalScale: number) => {
      const r = Math.max(3, Math.min(12, 3 + Math.sqrt(node.val ?? 1)));
      const isHighlighted = highlightNodes.size === 0 || highlightNodes.has(node.id);
      const colour = domainColour(node.domain);
      ctx.beginPath();
      ctx.arc(node.x!, node.y!, r, 0, 2 * Math.PI);
      ctx.fillStyle = isHighlighted ? colour : alpha(colour, 0.15);
      ctx.fill();
      ctx.strokeStyle = isHighlighted ? '#fff' : alpha('#fff', 0.15);
      ctx.lineWidth = 0.5;
      ctx.stroke();

      if (globalScale > 2.5 || highlightNodes.has(node.id)) {
        const fontSize = 10 / globalScale;
        ctx.font = `${fontSize}px Inter, sans-serif`;
        ctx.fillStyle = isHighlighted ? '#e2e8f0' : alpha('#e2e8f0', 0.25);
        ctx.textAlign = 'center';
        ctx.fillText(node.label.length > 24 ? node.label.slice(0, 22) + '…' : node.label, node.x!, node.y! + r + fontSize);
      }
    },
    [highlightNodes],
  );

  const linkColor = useCallback(
    (edge: GraphEdge) => {
      const isHighlighted = highlightEdges.size === 0 || highlightEdges.has(edge);
      const colour = relationColour(edge.relation);
      return isHighlighted ? colour : alpha(colour, 0.12);
    },
    [highlightEdges],
  );

  const linkWidth = useCallback(
    (edge: GraphEdge) => (highlightEdges.size === 0 || highlightEdges.has(edge) ? 1.5 : 0.5),
    [highlightEdges],
  );

  // ── Side panel ─────────────────────────────────────────────────────────────

  const sidePanel = selectedNode || selectedEdge ? (
    <Paper
      elevation={8}
      sx={{
        position: 'absolute',
        top: 16,
        right: 16,
        width: 320,
        maxHeight: 'calc(100% - 32px)',
        overflowY: 'auto',
        bgcolor: alpha('#0f172a', 0.95),
        border: '1px solid',
        borderColor: alpha('#94a3b8', 0.2),
        borderRadius: 2,
        p: 2,
        zIndex: 10,
      }}
    >
      <Stack direction="row" justifyContent="space-between" alignItems="flex-start">
        <Typography variant="subtitle2" color="primary.light" fontWeight={700}>
          {selectedNode ? 'Node' : 'Link'}
        </Typography>
        <IconButton size="small" onClick={clearHighlight} sx={{ color: 'text.secondary', mt: -0.5 }}>
          <CloseIcon fontSize="small" />
        </IconButton>
      </Stack>
      <Divider sx={{ my: 1, borderColor: alpha('#94a3b8', 0.15) }} />

      {selectedNode && (
        <Stack spacing={1}>
          <Typography variant="body2" color="text.primary" fontWeight={600} sx={{ wordBreak: 'break-word' }}>
            {selectedNode.label}
          </Typography>
          <Chip label={selectedNode.domain || 'unknown'} size="small"
            sx={{ bgcolor: alpha(domainColour(selectedNode.domain), 0.2), color: domainColour(selectedNode.domain), alignSelf: 'flex-start' }} />
          <Typography variant="caption" color="text.secondary">
            {selectedNode.val} connection{selectedNode.val !== 1 ? 's' : ''}
          </Typography>
        </Stack>
      )}

      {selectedEdge && (() => {
        const src = typeof selectedEdge.source === 'string' ? selectedEdge.source : (selectedEdge.source as GraphNode).label;
        const tgt = typeof selectedEdge.target === 'string' ? selectedEdge.target : (selectedEdge.target as GraphNode).label;
        return (
          <Stack spacing={1.5}>
            <Typography variant="body2" color="text.primary" sx={{ wordBreak: 'break-word' }}>
              <Box component="span" fontWeight={700}>{src}</Box>
            </Typography>
            <Chip
              label={selectedEdge.relation}
              size="small"
              sx={{ bgcolor: alpha(relationColour(selectedEdge.relation), 0.2), color: relationColour(selectedEdge.relation), alignSelf: 'flex-start' }}
            />
            <Typography variant="body2" color="text.primary" sx={{ wordBreak: 'break-word' }}>
              <Box component="span" fontWeight={700}>{tgt}</Box>
            </Typography>
            {selectedEdge.inverse_relation && (
              <Typography variant="caption" color="text.secondary">
                ← {selectedEdge.inverse_relation}
              </Typography>
            )}
            {selectedEdge.evidence_context && (
              <>
                <Divider sx={{ borderColor: alpha('#94a3b8', 0.15) }} />
                <Typography variant="caption" color="text.disabled" textTransform="uppercase" letterSpacing={0.8}>
                  Evidence context
                </Typography>
                <Typography variant="caption" color="text.secondary" sx={{ wordBreak: 'break-word' }}>
                  {selectedEdge.evidence_context}
                </Typography>
              </>
            )}
            {selectedEdge.source_title && (
              <>
                <Divider sx={{ borderColor: alpha('#94a3b8', 0.15) }} />
                <Typography variant="caption" color="text.secondary">
                  Source: {selectedEdge.source_title}
                </Typography>
              </>
            )}
            <Typography variant="caption" color="text.disabled">
              Confidence: {(selectedEdge.confidence * 100).toFixed(0)}%
            </Typography>
          </Stack>
        );
      })()}
    </Paper>
  ) : null;

  // ── Main render ────────────────────────────────────────────────────────────

  return (
    <Box sx={{ height: '100vh', display: 'flex', flexDirection: 'column', bgcolor: '#0a0f1e' }}>
      <SEO
        title="CHISG Knowledge Graph Explorer"
        description="Explore the CHISG semantic knowledge graph — bidirectional links between concepts, coloured by domain and relation type."
        path="/chisg/graph"
      />

      {/* Toolbar */}
      <Box sx={{ px: 2, py: 1.5, bgcolor: alpha('#0f172a', 0.95), borderBottom: '1px solid', borderColor: alpha('#94a3b8', 0.1), flexShrink: 0 }}>
        <Stack direction="row" alignItems="center" spacing={2} flexWrap="wrap" gap={1}>
          <Button component={Link} to="/CHISG" variant="text" size="small"
            sx={{ color: 'text.secondary', textTransform: 'none', minWidth: 0, p: 0.5 }}>
            ← CHISG
          </Button>
          <GraphIcon sx={{ color: 'primary.light', fontSize: 20 }} />
          <Typography variant="subtitle2" color="text.primary" fontWeight={700} sx={{ mr: 2 }}>
            Knowledge Graph Explorer
          </Typography>

          <FilterIcon sx={{ color: 'text.disabled', fontSize: 18 }} />

          <FormControl size="small" sx={{ minWidth: 130 }}>
            <InputLabel sx={{ color: 'text.secondary', fontSize: 13 }}>Domain</InputLabel>
            <Select value={domainFilter} label="Domain" onChange={e => setDomainFilter(e.target.value)}
              sx={{ color: 'text.primary', fontSize: 13, '.MuiOutlinedInput-notchedOutline': { borderColor: alpha('#94a3b8', 0.2) } }}>
              <MenuItem value=""><em>All domains</em></MenuItem>
              {availableDomains.map(d => <MenuItem key={d} value={d}>{d}</MenuItem>)}
            </Select>
          </FormControl>

          <FormControl size="small" sx={{ minWidth: 160 }}>
            <InputLabel sx={{ color: 'text.secondary', fontSize: 13 }}>Relation</InputLabel>
            <Select value={relationFilter} label="Relation" onChange={e => setRelationFilter(e.target.value)}
              sx={{ color: 'text.primary', fontSize: 13, '.MuiOutlinedInput-notchedOutline': { borderColor: alpha('#94a3b8', 0.2) } }}>
              <MenuItem value=""><em>All relations</em></MenuItem>
              {availableRelations.map(r => <MenuItem key={r} value={r}>{r}</MenuItem>)}
            </Select>
          </FormControl>

          <Tooltip title="Refresh">
            <IconButton size="small" onClick={fetchGraph} sx={{ color: 'text.secondary' }}>
              <RefreshIcon fontSize="small" />
            </IconButton>
          </Tooltip>

          {graphData && (
            <Typography variant="caption" color="text.disabled" sx={{ ml: 'auto' }}>
              {filteredGraphData.nodes.length} nodes · {filteredGraphData.links.length} edges
            </Typography>
          )}
        </Stack>
      </Box>

      {/* Graph canvas */}
      <Box ref={containerRef} sx={{ flex: 1, position: 'relative', overflow: 'hidden' }}>
        {loading && (
          <Box sx={{ position: 'absolute', inset: 0, display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: 2, zIndex: 20 }}>
            <CircularProgress size={40} thickness={3} />
            <Typography variant="body2" color="text.secondary">Loading graph…</Typography>
          </Box>
        )}

        {error && (
          <Box sx={{ position: 'absolute', inset: 0, display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 20 }}>
            <Typography color="error.main">{error}</Typography>
          </Box>
        )}

        {!loading && !error && graphData && (
          <ForceGraph2D
            ref={fgRef as any}
            width={dimensions.width}
            height={dimensions.height}
            graphData={filteredGraphData as any}
            backgroundColor="#0a0f1e"
            nodeId="id"
            linkSource="source"
            linkTarget="target"
            nodeCanvasObject={nodeCanvasObject as any}
            nodeCanvasObjectMode={() => 'replace'}
            onNodeClick={handleNodeClick as any}
            onLinkClick={handleEdgeClick as any}
            onBackgroundClick={clearHighlight}
            linkColor={linkColor as any}
            linkWidth={linkWidth as any}
            linkDirectionalArrowLength={4}
            linkDirectionalArrowRelPos={1}
            linkDirectionalArrowColor={linkColor as any}
            linkCurvature={0.15}
            nodePointerAreaPaint={(node: GraphNode, colour, ctx) => {
              const r = Math.max(5, Math.min(14, 4 + Math.sqrt(node.val ?? 1)));
              ctx.fillStyle = colour;
              ctx.beginPath();
              ctx.arc(node.x!, node.y!, r, 0, 2 * Math.PI);
              ctx.fill();
            }}
            cooldownTicks={80}
          />
        )}

        {sidePanel}

        {/* Legend */}
        <Box sx={{
          position: 'absolute', bottom: 16, left: 16,
          bgcolor: alpha('#0f172a', 0.85), border: '1px solid', borderColor: alpha('#94a3b8', 0.15),
          borderRadius: 1.5, p: 1.5, maxWidth: 200,
        }}>
          <Typography variant="caption" color="text.disabled" textTransform="uppercase" letterSpacing={0.8} display="block" mb={0.5}>
            Relation type
          </Typography>
          {['causes', 'inhibits', 'enables', 'is composed of', 'is used for', 'represents'].map(r => (
            <Stack key={r} direction="row" alignItems="center" spacing={0.75} sx={{ mb: 0.3 }}>
              <Box sx={{ width: 10, height: 3, bgcolor: relationColour(r), borderRadius: 1, flexShrink: 0 }} />
              <Typography variant="caption" color="text.secondary" noWrap>{r}</Typography>
            </Stack>
          ))}
          <Divider sx={{ my: 0.75, borderColor: alpha('#94a3b8', 0.1) }} />
          <Typography variant="caption" color="text.disabled" textTransform="uppercase" letterSpacing={0.8} display="block" mb={0.5}>
            Domain
          </Typography>
          {['microbiology', 'immunology', 'education', 'manual'].map(d => (
            <Stack key={d} direction="row" alignItems="center" spacing={0.75} sx={{ mb: 0.3 }}>
              <Box sx={{ width: 10, height: 10, bgcolor: domainColour(d), borderRadius: '50%', flexShrink: 0 }} />
              <Typography variant="caption" color="text.secondary" noWrap>{d}</Typography>
            </Stack>
          ))}
        </Box>
      </Box>
    </Box>
  );
};

export default GraphExplorer;
