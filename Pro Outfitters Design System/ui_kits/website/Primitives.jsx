/* Primitives.jsx — small, reusable building blocks for the Pro Outfitters website kit.
   Exposed on window so other Babel scripts can use them. */
const { useState, useEffect } = React;

const ArrowUR = ({ s = 16 }) => (
  <svg width={s} height={s} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6">
    <path d="M7 17 17 7M17 7H9M17 7v8" />
  </svg>
);

const Eyebrow = ({ children, tone }) => (
  <p className={"k-eyebrow" + (tone ? " " + tone : "")}>{children}</p>
);

const GoldRule = () => <hr className="k-gold-rule" />;

const Button = ({ variant = "primary", block, children, ...rest }) => (
  <button className={"k-btn " + variant + (block ? " block" : "")} {...rest}>{children}</button>
);

const LinkArrow = ({ children, onClick, small }) => (
  <a className="k-link" onClick={onClick}><span>{children}</span><ArrowUR s={small ? 14 : 16} /></a>
);

// Photographic placeholder. Real photography drops in here in production.
const Photo = ({ tone = "dusk", caption, style, className = "" }) => (
  <div className={"k-photo " + className} data-tone={tone} style={style}>
    {caption ? <span className="cap">{caption}</span> : null}
  </div>
);

const Wordmark = () => (
  <span className="k-wordmark">
    <span className="top">PRO&nbsp;OUTFITTERS</span>
    <span className="rule"><span className="ln"></span><span className="est">EST. 1970</span><span className="ln"></span></span>
  </span>
);

const SectionHead = ({ eyebrow, title, tone }) => (
  <div className="k-section-head">
    <GoldRule />
    <Eyebrow tone={tone}>{eyebrow}</Eyebrow>
    <h2>{title}</h2>
  </div>
);

Object.assign(window, { ArrowUR, Eyebrow, GoldRule, Button, LinkArrow, Photo, Wordmark, SectionHead });
