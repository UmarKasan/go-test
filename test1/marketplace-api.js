// marketplace-api.js

const API_BASE_URL = 'http://localhost:8080/api';

class MarketplaceAPI {
    static async getAllItems() {
        try {
            const response = await fetch(`${API_BASE_URL}/items`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return await response.json();
        } catch (error) {
            console.error("Could not get items:", error);
        }
    }

    static async getItem(id) {
        try {
            const response = await fetch(`${API_BASE_URL}/items/${id}`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return await response.json();
        } catch (error) {
            console.error(`Could not get item ${id}:`, error);
        }
    }

    static async createItem(item) {
        try {
            const response = await fetch(`${API_BASE_URL}/items`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(item),
            });
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return await response.json();
        } catch (error) {
            console.error("Could not create item:", error);
        }
    }

    static async updateItem(id, item) {
        try {
            const response = await fetch(`${API_BASE_URL}/items/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(item),
            });
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return await response.json();
        } catch (error) {
            console.error(`Could not update item ${id}:`, error);
        }
    }

    static async deleteItem(id) {
        try {
            const response = await fetch(`${API_BASE_URL}/items/${id}`, {
                method: 'DELETE',
            });
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return true;
        } catch (error) {
            console.error(`Could not delete item ${id}:`, error);
        }
    }
}

// Example usage:
async function exampleUsage() {
    // Get all items
    const allItems = await MarketplaceAPI.getAllItems();
    console.log("All items:", allItems);

    // Create a new item
    const newItem = {
        FirstName: "John",
        LastName: "Doe",
        Product: "Recycled Paper",
        Quantity: 100,
        Condition: "New",
        CollectionLocation: "789 Green Street"
    };
    const createdItem = await MarketplaceAPI.createItem(newItem);
    console.log("Created item:", createdItem);

    // Get a specific item
    if (createdItem && createdItem.id) {
        const retrievedItem = await MarketplaceAPI.getItem(createdItem.id);
        console.log("Retrieved item:", retrievedItem);

        // Update the item
        retrievedItem.Quantity = 150;
        const updatedItem = await MarketplaceAPI.updateItem(retrievedItem.id, retrievedItem);
        console.log("Updated item:", updatedItem);

        // Delete the item
        const deleted = await MarketplaceAPI.deleteItem(retrievedItem.id);
        console.log("Item deleted:", deleted);
    }
}

// Uncomment the line below to run the example usage
// exampleUsage();