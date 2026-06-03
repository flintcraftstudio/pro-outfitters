/* InquiryModal.jsx — the one genuinely interactive flow in the kit:
   a 3-step "Plan Your Trip" inquiry, fake-submitted to a success state. */
const { useState: useS3 } = React;

const CloseIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5"><path d="M18 6 6 18M6 6l12 12" /></svg>
);
const Check = () => (
  <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5"><path d="M20 6 9 17l-5-5" /></svg>
);

const InquiryModal = ({ onClose }) => {
  const [step, setStep] = useS3(0);
  const [pursuit, setPursuit] = useS3("Fly Fishing");
  const [done, setDone] = useS3(false);

  const next = () => setStep((s) => Math.min(2, s + 1));
  const back = () => setStep((s) => Math.max(0, s - 1));

  return (
    <div className="k-overlay" onClick={(e) => { if (e.target === e.currentTarget) onClose(); }}>
      <div className="k-modal" role="dialog" aria-modal="true">
        <button className="close" onClick={onClose} aria-label="Close"><CloseIcon /></button>

        {done ? (
          <div className="k-success">
            <div className="mark"><Check /></div>
            <h3 style={{ marginTop: 0 }}>We'll be in touch.</h3>
            <p className="sub" style={{ marginBottom: "1.6rem" }}>
              Thank you. One of the family will reach out within two business days to plan the details of your trip.
            </p>
            <Button variant="secondary" block onClick={onClose}>Close</Button>
          </div>
        ) : (
          <React.Fragment>
            <Eyebrow tone="gold">Plan Your Trip · {step + 1} of 3</Eyebrow>

            {step === 0 && (<React.Fragment>
              <h3>What are you after?</h3>
              <p className="sub">There's no wrong answer — many guests do both.</p>
              <div className="k-choice" style={{ marginBottom: "1.6rem" }}>
                {["Fly Fishing", "Bird Hunting", "Both"].map((p) => (
                  <button key={p} className={pursuit === p ? "on" : ""} onClick={() => setPursuit(p)}>{p}</button>
                ))}
              </div>
              <Button variant="primary" block onClick={next}>Continue</Button>
            </React.Fragment>)}

            {step === 1 && (<React.Fragment>
              <h3>When &amp; how many?</h3>
              <p className="sub">Approximate is fine. We'll match you to the right water and lodge.</p>
              <div className="k-field"><label>Preferred dates</label><input type="text" placeholder="e.g. late September, 4 nights" /></div>
              <div className="k-field"><label>Party size</label>
                <select defaultValue="2"><option>1</option><option>2</option><option>3–4</option><option>5–6</option><option>7+</option></select>
              </div>
              <div style={{ display: "flex", gap: ".8rem" }}>
                <Button variant="secondary" onClick={back}>Back</Button>
                <Button variant="primary" block onClick={next}>Continue</Button>
              </div>
            </React.Fragment>)}

            {step === 2 && (<React.Fragment>
              <h3>Where do we reach you?</h3>
              <p className="sub">A real person reads every inquiry — no automated funnels.</p>
              <div className="k-field"><label>Name</label><input type="text" placeholder="Your name" /></div>
              <div className="k-field"><label>Email</label><input type="email" placeholder="you@email.com" /></div>
              <div className="k-field"><label>Anything we should know?</label><textarea rows="2" placeholder="Optional"></textarea></div>
              <div style={{ display: "flex", gap: ".8rem" }}>
                <Button variant="secondary" onClick={back}>Back</Button>
                <Button variant="primary" block onClick={() => setDone(true)}>Send Inquiry</Button>
              </div>
            </React.Fragment>)}
          </React.Fragment>
        )}
      </div>
    </div>
  );
};

Object.assign(window, { InquiryModal });
