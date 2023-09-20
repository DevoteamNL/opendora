import React, {  useEffect } from 'react';
import {  Grid } from '@material-ui/core';
import {
  Header,
  Page,
  Content,
  ContentHeader,
  SupportButton,
} from '@backstage/core-components';
import './DashboardComponent.css';
import { HighlightTextBoxComponent } from '../HighlightTextBoxComponent/HighlightTextBoxComponent';
import SimpleCharts from '../BarChartComponent/BarChartComponent';
import DropdownComponent from '../DropdownComponent/DropdownComponent';
import { useParams } from 'react-router';

import GroupDataService from '../../services/GroupDataService';
import { DeploymentFrequencyData } from '../../models/DeploymentFrequencyData';
import { calculateAverage} from '../../utils/helpers';

export  const  DashboardComponent = () => {
  const [chartData, setChartData] = React.useState<DeploymentFrequencyData | null>(null);


  const params= useParams();
  const componentName = params.name
  const [groupQueryParam, setGroupQueryParam] = React.useState<string | null>(null);
  const [selectedTimeUnit, setSelectedTimeUnit] = React.useState('Weekly');
  const [deploymentAverage, setDeploymentAverage] = React.useState(0);


  useEffect(() => {
    
    // fetching group
    GroupDataService.getAncestry(componentName).then((res) => {
      setGroupQueryParam(res.relations.filter((relation: any) => { return relation.type == "ownedBy"}).map((filteredRelation:any)=> {return filteredRelation.target.name})[0])
    })
  }, []);


  useEffect(() => {
      if(groupQueryParam)
      GroupDataService.getMockData(groupQueryParam, selectedTimeUnit).then((response: DeploymentFrequencyData) => {
        // here we get fetch data for the graphs
        setChartData(response);
        setDeploymentAverage(calculateAverage(response));

     })
  }, [groupQueryParam, selectedTimeUnit])

  return (
    <Page themeId="tool">
      <Header title="Devoteam DORA plugin" subtitle="Through insight to perfection">
        {/* <HeaderLabel label="Owner" value="Team X" />
        <HeaderLabel label="Lifecycle" value="Alpha" /> */}
      </Header>
      <Content>
        <ContentHeader title="DORA metrics">
          <SupportButton>Pluging for displaying DORA metrics</SupportButton>
        </ContentHeader>
        <Grid container spacing={3} direction="column">
          <Grid container>
            <Grid item xs={12} className="gridBorder">
              <div className="gridBoxText">
                <h1>Deployment statistics</h1>
                <p>Analysis of successful deployments and CFR</p>
                <Grid container>
                  <Grid item xs={4}>
                    <DropdownComponent
                      onSelect={setSelectedTimeUnit}
                      selection={selectedTimeUnit}
                      type="timeUnit"
                    />
                  </Grid>
                </Grid>
              </div>
            </Grid>


            <Grid item xs={6} className="gridBorder">
              <div className="gridBoxText">
                <HighlightTextBoxComponent
                  title={selectedTimeUnit + " average number of deployments"}
                  highlight={deploymentAverage}
                  textColour="positiveHighlight"
                />
              </div>
            </Grid>
            <Grid item xs={6} className="gridBorder">
              <div className="gridBoxText">
                {/* <p>Overall change failure rate</p>
                    <h1> 3%</h1> */}
                <HighlightTextBoxComponent
                  title="Overall change failure rate"
                  highlight="5.2%"
                  text="*calculated on failures and incidents"
                  textColour="warning"
                />
              </div>
            </Grid>

            <Grid item xs={12} className="gridBorder">
              <div className="gridBoxText">
                <SimpleCharts deploymentFrequencyData={chartData} />
              </div>
            </Grid>
          </Grid>

          <Grid item></Grid>
        </Grid>
      </Content>
    </Page>
  );
};
