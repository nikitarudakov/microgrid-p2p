export const useEnergyResourcesStore= defineStore(
    "energyResourcesStore",
    () => {
      const energyResources = useState<EnergyResource[]>()

      async function fetch() {
          const { data} = await useAsyncGql({ operation: 'EnergyResources' });
          energyResources.value = data.value.energyResources
      }

      function updateEnergyResource(updatedResource: EnergyResource) {
          energyResources.value = energyResources.value.map((resource) =>
              resource.id === updatedResource.id ? updatedResource : resource
          );
      }

      return {energyResources, fetch, updateEnergyResource};
    },
)

export interface EnergyResource {
    id: string,
    name: string,
    capacity: number,
    price: number,
    producer: User
}

export interface User {
    first_name: string,
    last_name: string
}
