import { gql } from "@apollo/client/core"

export const ENERGY_RESOURCES = gql`
  query EnergyResources {
    energyResources {
      id
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

export const PURCHASE_REQUEST = gql`
  mutation PurchaseEnergy ($id: ID!, $capacity: Float!) {
    purchaseEnergy(in: {
      id: $id
      capacity: $capacity
    }) {
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
