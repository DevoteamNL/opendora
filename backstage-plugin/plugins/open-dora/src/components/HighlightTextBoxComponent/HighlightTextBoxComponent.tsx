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
      <h1 className="margin-left-offset-m25 headerStyle">{props.title}</h1>
      <div className={props.healthStatus}>
        <div className="highlight">{props.highlight}</div>
      </div>
      <div className="notification">{props.text}</div>
    </div>
  );
};
