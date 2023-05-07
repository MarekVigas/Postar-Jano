import "little-state-machine";

declare module "little-state-machine" {
  interface GlobalState {
    promo: string | null
  }
}
