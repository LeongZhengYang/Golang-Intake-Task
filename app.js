const runningNumberElement = document.getElementById('runningNumber');
const startButton = document.getElementById('startButton');
const inputNumberElement = document.getElementById('inputNumber');
const resultElement = document.getElementById('result');
const accumulatedValueElement = document.getElementById('accumulatedValue');

let randomNumberInterval;
let collectedDigits = [];
let collectionActive = false;

function generateRandomNumber() {
    return Math.floor(Math.random() * 9000) + 1000;
}

function updateRunningNumber() {
    if (collectionActive) {
        const randomNumber = generateRandomNumber();
        runningNumberElement.textContent = randomNumber.toString();
        console.log("Generated number:", randomNumber); 
    }
}

function startDataCollection() {
    const inputNumber = inputNumberElement.value;

    if (!inputNumber || inputNumber.length !== 5 || isNaN(inputNumber)) {
        alert("Please enter a 5-digit number.");
        return;
    }

    collectedDigits = [];
    collectionActive = true;
    let count = 0;

    // Reset the result display
    resultElement.textContent = '';
    
    // Update accumulated value display
    accumulatedValueElement.textContent = `Accumulated Value: `;

    // Clear any existing interval
    if (randomNumberInterval) {
        clearInterval(randomNumberInterval);
    }

    // Start new interval for updating running number
    randomNumberInterval = setInterval(updateRunningNumber, 1000);

    const collectionInterval = setInterval(() => {
        if (count < 5) {
            const lastDigit = runningNumberElement.textContent.slice(-1);
            if (lastDigit) { 
                collectedDigits.push(lastDigit);
                accumulatedValueElement.textContent = `Accumulated Value: ${collectedDigits.join('')}`;
                count++;
            }
        } else if (count == 5) {
            clearInterval(collectionInterval);
            clearInterval(randomNumberInterval);
            collectionActive = false;
            storeTransaction();
        }
    }, 60000); 
}

function storeTransaction() {
    const inputNumber = inputNumberElement.value;

    if (inputNumber.length !== 5) {
        alert("Please enter a 5-digit number.");
        return;
    }

    const accumulatedNumber = collectedDigits.join('');
    
    // Send to the server (Golang backend) using Fetch API
    fetch('/store', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            inputNumber: inputNumber,
            accumulatedNumber: accumulatedNumber,
        }),
    })
    .then(response => response.json())
    .then(data => {
        resultElement.textContent = `Matched: ${data.matched}, Continuous: ${data.continuous}, Permutation: ${data.permutation}`;
    })
    .catch(error => {
        console.error('Error:', error);
    });
}

startButton.addEventListener('click', () => {
    if (!collectionActive) {
        startDataCollection();
    }
});