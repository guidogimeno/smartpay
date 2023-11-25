# SmartPay

## Description

This is a backend application written in Go that provides a convenient way to calculate whether it is more advantageous to pay an amount in installments or in a lump sum. By considering the inflation index and the Annual Nominal Rate (TNA) provided by the bank, this application helps users make informed financial decisions.

## Features

- Calculates the optimal payment method (lump sum or installments) based on the given parameters: amount, number of installments, and interest rate.
- Retrieves the current inflation index and TNA from the bank.
- Performs the necessary calculations to determine the most cost-effective payment option.
- Provides clear and concise results to assist users in their financial decision-making process.

## Technologies Used

- Go programming language
- HTTP client for API communication
- BCRA's inflation index and TNA data source

## Usage

1. Ensure you have Go installed on your system.
2. Clone the repository and navigate to the project directory.
3. Configure the necessary bank API credentials in the application.
4. Build and run the application.
5. Access the application's endpoints to perform calculations based on your requirements.
6. Review the results and make informed decisions regarding your payments.

## Future Enhancements

- Implement a user-friendly interface for easy interaction.
- Add support for multiple banks' data sources.
- Enhance error handling and validation for input parameters.
- Implement caching mechanisms to improve performance.
- Incorporate historical data analysis for more accurate calculations.

## Contributors

- [Guido Gimeno](https://github.com/guidogimeno) - Project Lead

Please feel free to contribute by submitting bug reports, feature requests, or pull requests to this project.
