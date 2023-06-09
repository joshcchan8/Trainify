import React from "react"
import axios from "axios"

function Item(props) {
    return (
        <tr>
            <td>{props.item.item_name}</td>
            <td>{props.item.difficulty}</td>
            <td>{props.item.minutes}</td>
            <td>{props.item.calories_burned}</td>
            <td>{props.item.targeted_muscle_groups}</td>
            <td>{props.item.workout_description}</td>
        </tr>
    )
}

export default function Dashboard( props ) {

    const [items, setItems] = React.useState([])

    React.useEffect(() => {
        setItems(props.items)
    }, [props.items])

    function itemList() {
        return items.map(currentItem => {
            return (
                <Item
                    item={currentItem}
                    key={currentItem.item_id}
                />
            )
        })
    }

    return (
        <div>
            <h3>Item List</h3>
            <table>
                <thead>
                    <tr>
                        <th>Item Name</th>
                        <th>Difficulty</th>
                        <th>Minutes</th>
                        <th>Calories Burned</th>
                        <th>Targeted Muscle Groups</th>
                        <th>Workout Description</th>
                    </tr>
                </thead>
                <tbody>
                    { itemList() }
                </tbody>
            </table>
        </div>
    )
}