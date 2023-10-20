import React from 'react';
import './HighlightTextBoxComponent.css';

interface HighlightTextBoxComponentProps {
  title: string;
  textColour: 'warning' | 'critical' | 'positiveHighlight';
  highlight: string;
  text?: string;
}
export const HighlightTextBoxComponent = (
  props: HighlightTextBoxComponentProps,
) => {
  return (
    <div>
      <h1>{props.title}</h1>
      <div className={props.textColour}>
        <div className="highlight">{props.highlight}</div>
      </div>
      <div className="notification">{props.text}</div>
    </div>
  );
};
