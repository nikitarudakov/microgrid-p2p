import { gql } from "@apollo/client/core"

export const GET_ENERGY_RESOURCES = gql`
  query GetEnergyResources($ownerName: String!) {
    energy_resources(owner_name: $ownerName){
        id
        owner_name
        capacity
    }
  }
`;