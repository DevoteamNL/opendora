import React from 'react'
import './HighlightTextBoxComponent.css'

export const HighlightTextBoxComponent = (prop: any) => {
    
    // console.log("title", title.title, title.highlight, text)

    if (prop.warning) {

    }
    
    return(
    <div >
        <h1>{prop.title}</h1>
        <div className={prop.textColour}>
            <div className="highlight">{prop.highlight}</div>
        </div>
        <div className="notification">{prop.text}</div>  
    </div>
)}