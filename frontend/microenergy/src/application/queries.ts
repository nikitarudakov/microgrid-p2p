import { gql } from "@apollo/client/core"

export const ENERGY_RESOURCES = gql`
  query EnergyResources {
    energyResources {
      name
      producer {
        first_name
        last_name
      }
      capacity
      price
    }
  }
`;