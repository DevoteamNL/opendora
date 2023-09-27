import { DeploymentFrequencyData } from "../models/DeploymentFrequencyData";

export function calculateAverage(data: DeploymentFrequencyData | null) {
    if (data === null){
      return "0"
    }
    const totalDeployments = data.dataPoints.reduce((a,v) => a + v.value, 0);
    return Number((totalDeployments/data.dataPoints.length).toFixed(2));
  }