import React from 'react';
import { Helmet } from 'react-helmet-async';

interface SEOProps {
  title: string;
  description: string;
  path: string;
  /** Override the default OG image (relative to domain root, e.g. /ESPLogoLong.png) */
  image?: string;
  /** Additional structured data (JSON-LD) */
  jsonLd?: Record<string, unknown>;
}

const DOMAIN = 'https://espthinking.co.uk';
const DEFAULT_IMAGE = '/ESPLogoLong.png';

/**
 * Per-page SEO tags using react-helmet-async.
 * Overrides the defaults set in public/index.html.
 */
const SEO: React.FC<SEOProps> = ({ title, description, path, image, jsonLd }) => {
  const fullTitle = `${title} | ESP Thinking`;
  const url = `${DOMAIN}${path}`;
  const imageUrl = `${DOMAIN}${image || DEFAULT_IMAGE}`;

  return (
    <Helmet>
      <title>{fullTitle}</title>
      <meta name="description" content={description} />
      <link rel="canonical" href={url} />

      {/* Open Graph */}
      <meta property="og:title" content={fullTitle} />
      <meta property="og:description" content={description} />
      <meta property="og:url" content={url} />
      <meta property="og:image" content={imageUrl} />
      <meta property="og:type" content="website" />
      <meta property="og:site_name" content="ESP Thinking" />

      {/* Twitter */}
      <meta name="twitter:card" content="summary_large_image" />
      <meta name="twitter:title" content={fullTitle} />
      <meta name="twitter:description" content={description} />
      <meta name="twitter:image" content={imageUrl} />

      {/* JSON-LD structured data */}
      {jsonLd && (
        <script type="application/ld+json">
          {JSON.stringify(jsonLd)}
        </script>
      )}
    </Helmet>
  );
};

export default SEO;
