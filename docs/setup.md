# Setting up

1. Using this app requires you have a Xero account and Xero organisation.

2. To track your spending and income on Xero, you should connect your personal bank account via [Bank feeds](https://central.xero.com/s/article/Bank-feeds).

3. Categorise your spending and income into appropriate accounts with "Expenses" and "Revenue" account types.
   For example, you can create a "Groceries" account, with the "Expenses" type and "Salary" account with the "Revenue" type.

4. Set budgets for each your accounts under Accounting > Reports > Budget Manager

5. To connect tui-budgeter to you must [create a Xero app](https://developer.xero.com/myapps/). The app uses OAuth2 PKCE flow.

6. Create a `.env` file following `.env.example`. Set `XERO_CLIENT_ID` to the client id of the app you created above.

7. tui-budgeter should now work. You can now track your spending and income right from your terminal! Hooray!
