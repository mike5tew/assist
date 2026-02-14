import React from 'react';

type Props = {
  width?: number | string;
  height?: number | string;
  color?: string;
  title?: string;
};

const DivideIcon: React.FC<Props> = ({ width = 36, height = 36, color = 'currentColor', title = 'Divide' }) => (
  <svg
    width={width}
    height={height}
    viewBox="-6 -6 60 60"
    preserveAspectRatio="xMidYMid meet"
    style={{ overflow: 'visible', display: 'block' }}
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    role="img"
    aria-label={title}
  >
    <title>{title}</title>
    {/* Left person */}
    <circle cx="11" cy="12" r="3.6" stroke={color} strokeWidth="2.2" />
    <path d="M8 22c0-3 6-3 6 0v6" stroke={color} strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round" />
    <path d="M6 18c3 2 7 2 10 0" stroke={color} strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round" />

    {/* Right person */}
    <circle cx="37" cy="12" r="3.6" stroke={color} strokeWidth="2.2" />
    <path d="M40 22c0-3-6-3-6 0v6" stroke={color} strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round" />
    <path d="M42 18c-3 2-7 2-10 0" stroke={color} strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round" />

    {/* Jagged divide */}
    <path d="M24 6 L28 14 L22 18 L26 26 L20 30 L28 42" stroke={color} strokeWidth="2.6" strokeLinecap="round" strokeLinejoin="round" fill="none" />
  </svg>
);

export default DivideIcon;
