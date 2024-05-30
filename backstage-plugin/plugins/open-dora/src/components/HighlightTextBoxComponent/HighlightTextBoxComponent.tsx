import React from 'react';
import './HighlightTextBoxComponent.css';

interface HighlightTextBoxComponentProps {
  title: string;
  healthStatus: string;
  highlight: string;
  text?: string;
}
export const HighlightTextBoxComponent = (
  props: HighlightTextBoxComponentProps,
) => {
  return (
    <div className="highlightTextBoxBorder">
      <div className={props.healthStatus}>
        <div className="highlight">{props.highlight}</div>
      </div>
      <span className="headerStyle">{props.title}</span>
      <div className="notification">{props.text}</div>
    </div>
  );
};
