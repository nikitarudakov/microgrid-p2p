export interface User {
  firstName: string;
  lastName: string;
}

export interface EnergyResource {
  id: string,
  name: string,
  producer: User,
  capacity: number,
  price: number
}
