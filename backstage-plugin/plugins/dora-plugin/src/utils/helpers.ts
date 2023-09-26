import { DeploymentFrequencyData } from "../models/DeploymentFrequencyData";

export function calculateAverage(data: DeploymentFrequencyData) {
    const totalDeployments = (data.dataPoints.reduce((a,v) => a = a + v.value, 0))
    return Number((totalDeployments/data.dataPoints.length).toFixed(2));
  }