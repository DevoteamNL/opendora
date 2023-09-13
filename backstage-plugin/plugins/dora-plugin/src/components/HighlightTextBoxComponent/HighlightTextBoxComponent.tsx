import React from 'react';
import './HighlightTextBoxComponent.css';

export const HighlightTextBoxComponent = (prop: any) => {
  if (prop.warning) {
  }

  return (
    <div>
      <h1>{prop.title}</h1>
      <div className={prop.textColour}>
        <div className="highlight">{prop.highlight}</div>
      </div>
      <div className="notification">{prop.text}</div>
    </div>
  );
};
