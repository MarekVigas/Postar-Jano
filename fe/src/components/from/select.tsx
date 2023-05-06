import { IonSelect, IonSelectOption } from "@ionic/react";
import { Control, Controller, UseFormSetValue } from "react-hook-form";
import { Registration } from "../../utils/types";

type SelectOptions = {
  label: string,
  value: string | boolean
}

type RadioProps = {
  control: Control<Registration>
  setValue: UseFormSetValue<Registration>
  name: any
  options: SelectOptions[]
  placeholder: string
  okText: string
  cancelText: string
}

export const Select: React.FC<RadioProps> = ({ control, setValue, name, options, placeholder, okText, cancelText }) => (
  <Controller
    render={({ field }) => (
      <IonSelect
        placeholder={placeholder}
        okText={okText}
        cancelText={cancelText}
        value={field.value}
        onIonChange={e => setValue(name, e.detail.value)}
      >
        {options.map(({ value, label }, index) => (
          <IonSelectOption key={index} value={value}>{label}</IonSelectOption>
        ))}
      </IonSelect>
    )}
    control={control}
    name={name}
  />
)