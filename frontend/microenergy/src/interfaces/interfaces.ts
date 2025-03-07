export interface User {
  firstName: string;
  lastName: string;
}

export interface EnergyResource {
  name: string,
  producer: User,
  capacity: number,
  price: number
}
