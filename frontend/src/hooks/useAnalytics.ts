import { useEffect, useCallback, useRef } from 'react';
import { useLocation } from 'react-router-dom';

const API_BASE = '/api/analytics';

// Fire-and-forget beacon — never blocks the UI
const beacon = (endpoint: string, body: object) => {
  try {
    const blob = new Blob([JSON.stringify(body)], { type: 'application/json' });
    if (navigator.sendBeacon) {
      navigator.sendBeacon(`${API_BASE}/${endpoint}`, blob);
    } else {
      fetch(`${API_BASE}/${endpoint}`, {
        method: 'POST',
        body: JSON.stringify(body),
        headers: { 'Content-Type': 'application/json' },
        keepalive: true,
      }).catch(() => {});
    }
  } catch {
    // Analytics should never throw
  }
};

/**
 * useAnalytics — drop into App.tsx to auto-track page views
 * and get a trackEvent() function for CTA clicks.
 *
 * Usage:
 *   const { trackEvent } = useAnalytics();
 *   trackEvent('contact_open', { label: 'PrimaryOS CTA' });
 */
export function useAnalytics() {
  const location = useLocation();
  const prevPath = useRef<string | null>(null);

  // Auto-track page views on route change
  useEffect(() => {
    const path = location.pathname;

    // Avoid duplicate fires for the same path
    if (path === prevPath.current) return;
    prevPath.current = path;

    beacon('pageview', {
      path,
      referrer: document.referrer || '',
      title: document.title,
    });
  }, [location.pathname]);

  // Manual event tracking
  const trackEvent = useCallback(
    (name: string, meta?: Record<string, string>) => {
      beacon('event', {
        name,
        path: location.pathname,
        meta,
      });
    },
    [location.pathname],
  );

  return { trackEvent };
}

/**
 * Standalone function for tracking events outside React components
 * (e.g., in utility functions).
 */
export function trackEvent(name: string, path: string, meta?: Record<string, string>) {
  beacon('event', { name, path, meta });
}

export default useAnalytics;
