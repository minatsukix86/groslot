
# 🔍 Discord Token Checker

Un script en Go pour vérifier la validité de tokens Discord et afficher les informations associées à chaque compte, y compris le type de Nitro et les amis en ligne.

## ⚙️ Fonctionnalités

- **Vérification de validité de token** 🔑
- **Identification du type de Nitro** 🎉 (Basic, Boost, ou aucun)
- **Compteur d'amis et statut en ligne** 👥
- **Pause aléatoire entre les vérifications** ⏲️ 

## 🛠️ Utilisation

1. Placez vos tokens dans un fichier texte (ex. `tokens.txt`), un token par ligne.
2. Exécutez le script avec :
   ```bash
   ./discord-check.exe tokens.txt
   ```
3. Consultez les résultats et la séparation entre chaque compte.

## 📜 Exemple de sortie

```
Token valide pour l'utilisateur : JohnDoe#1234
L'utilisateur possède un Nitro Boost.
Total d'amis : 200
-----------------------------------------------------------
```

---

Ce script est utile pour les développeurs qui veulent rapidement vérifier des informations de compte Discord.
