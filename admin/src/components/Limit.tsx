import React from 'react'
import {IDay, IStat} from "../api/events";
import {ProgressBar} from "react-bootstrap";

interface Props {
    day :IDay;
    stat :IStat;
}

const Limit:React.FC<Props> = ({day, stat}) => {
    const splitLimits = day.limit_boys != null || day.limit_girls != null

    const variantForPercentage = (now :number) =>  {
        if (now < 50) return "info"
        if (now < 80) return "success"
        if (now < 100) return "warning"
        return "danger"
    }

    if (splitLimits) {
        return (
            <div>{stat.boys_count}CH + {stat.girls_count}D / {stat.capacity}(max.{stat.limit_boys}CH + max.{stat.limit_girls}D)</div>
        )
    }
    const percentage = (stat.boys_count + stat.girls_count) === 0 ? 0 : (stat.boys_count + stat.girls_count)/stat.capacity * 100

    return (
        <div>
            {stat.boys_count + stat.girls_count}/{stat.capacity}
            <ProgressBar now={percentage} variant={variantForPercentage(percentage)}/>
        </div>
    )

}

export default Limit
