import React, { useState, useEffect } from 'react';
import { Typography, Grid } from '@material-ui/core';
import {
  InfoCard,
  Header,
  Page,
  Content,
  ContentHeader,
  HeaderLabel,
  SupportButton,
} from '@backstage/core-components';
import "./DashboardComponent.css"
import { HighlightTextBoxComponent } from '../HighlightTextBoxComponent/HighlightTextBoxComponent';
import SimpleCharts from '../BarChartComponent/BarChartComponent';
import DropdownComponent from '../DropdownComponent/DropdownComponent';

import GroupDataService from '../../services/GroupDataService'



export const DashboardComponent = () => {

  const [chartData, setChartData] = React.useState('')

  

  useEffect(()=> {
    GroupDataService.getMockData()
      .then((response: any) => { 
        // here we get fetch data for the graphs
        setChartData(response);
  });
    }, [])

  const [selectedGroup, setSelectedGroup] = React.useState('');
  const updateGrp = (grp) => {
    setSelectedGroup(grp);
    // TODO
    // here we refresh query based on grp
  }

  const [selectedTimeUnit, setSelectedTimeUnit] = React.useState('');
  const updateTimeUnit = (timeUnit) => {
    setSelectedTimeUnit(timeUnit);
    // TODO
    // here we refresh query based on grp
  }




  return (
    <Page themeId="tool">
      <Header title="Welcome to dora-plugin!" subtitle="Optional subtitle">
        <HeaderLabel label="Owner" value="Team X" />
        <HeaderLabel label="Lifecycle" value="Alpha" />
      </Header>
      <Content>
        <ContentHeader title="Plugin title">
          <SupportButton>A description of your plugin goes here.</SupportButton>
        </ContentHeader>
        <Grid container spacing={3} direction="column">
          
          <Grid container>
              <Grid item xs={12} className="gridBorder">
                <div className="gridBoxText" >
                  <h1>Deployment statistics</h1>
                  <p>Analysis of successful deployments and CFR</p>   
                  <Grid container>
                    <Grid item xs={6}>
                     <DropdownComponent onSelect={(e)=> {updateGrp(e)}} selection = {selectedGroup} type = "group"/>
                    </Grid>
                    <Grid item xs={6}>
                     <DropdownComponent onSelect={(e)=> {updateTimeUnit(e)}} selection = {selectedTimeUnit} type = "timeUnit"/>
                    </Grid>
                  </Grid>
                 
                </div>

              </Grid>


              <Grid item xs={4} className="gridBorder">
              </Grid>

              <Grid item xs={4} className="gridBorder">
                  <div className="gridBoxText" >
                  <HighlightTextBoxComponent title="Average number of deployments per week" highlight="31" textColour='positiveHighlight'/>
                  </div>
                
              </Grid>
              <Grid item xs={4} className="gridBorder">
                <div className="gridBoxText" >
                    {/* <p>Overall change failure rate</p>
                    <h1> 3%</h1> */}
                    <HighlightTextBoxComponent title="Overall change failure rate" highlight="5.2%" text="*calculated on failures and incidents" textColour='warning'/>
                </div>
                
              </Grid>

              <Grid item xs={6} className="gridBorder">
                <div className="gridBoxText" >
                 <SimpleCharts ChartData={chartData}/>
                </div>
              </Grid>

          </Grid>


          <Grid item>
          </Grid>
        </Grid>
      </Content>
    </Page>)
}