import React from "react";

interface Props {
   fields :IViewFields;
   setFields :React.Dispatch<React.SetStateAction<IViewFields>>
}

export interface IViewFields {
    [key: string]: ISetting;
}

interface ISetting {
    display :string;
    show :boolean;
}


const ViewFilter: React.FC<Props> = ({fields, setFields}) => {
    const handleChange = (name:string, checked :boolean) => {
        setFields(prevState => {
                const updated = {...prevState}
                updated[name] = {...prevState[name], show:checked}
                return updated
            }
        )
    }

    return (
        <div className="filters">
            {Object.keys(fields).map(
                (k,index) => <div key={index}>
                        <input
                            type="checkbox"
                            checked={fields[k].show}
                            onChange={(e)=>{
                                handleChange(k, e.target.checked)
                            }}
                        />
                        <span>{fields[k].display}</span>
                    </div>

                )
            }
        </div>
    )
}
export default ViewFilter