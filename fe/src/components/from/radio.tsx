import { IonRadio, IonRadioGroup } from "@ionic/react";
import { UseFormGetValues, UseFormRegister, UseFormSetValue } from "react-hook-form";
import { Registration } from "../../utils/types";

type RadioOption = {
  label: string,
  value: string | boolean
}

type RadioProps = {
  register: UseFormRegister<Registration>
  setValue: UseFormSetValue<Registration>
  getValues: UseFormGetValues<Registration>
  name: any
  options: RadioOption[]
}

export const RadioGroup: React.FC<RadioProps> = ({ register, setValue, getValues, name, options }) => (
  <IonRadioGroup
    style={{ display: 'flex', flexDirection: 'column', width: '100%', justifyContent: 'space-between', height: '6em' }}
    {...register(name)}
    defaultValue={getValues(name)}
    value={getValues(name)}
    onIonChange={e => setValue(name, e.detail.value)}
  >
    {
      options.map(({ value, label }, index) => (
        <IonRadio key={index} slot="start" labelPlacement='end' justify='start' value={value}>{label}</IonRadio>
      ))
    }
  </IonRadioGroup>
)