# My Body Tracker

Réalisé par [Joris VILARDELL](https://github.com/ZUHOWKS) & [Damien RODRIGUEZ](https://github.com/rodriguezdamien)

**My Body Tracker** est une application de suivi nutritionnel en Go permettant de calculer l'IMC, suivre les dépenses énergétiques et adapter les besoins caloriques en fonction du niveau d'activité physique.

Cette solution s’adresse à un public soucieux de sa performance physique ou de son équilibre
alimentaire, qui a besoin d’un outil fiable, rapide et sans fonctionnalités superflues. Il est envisageable
que l’application, à terme, puisse ensuite être intégrée à une interface mobile ou web.

## Utilisation

Pour lancer le projet **My Body Tracker** développé en Go, voici une procédure détaillée :

---

### 1. **Cloner le dépôt Git**

Commence par cloner le dépôt depuis GitHub :

```bash
git clone https://github.com/ZUHOWKS/my-body-tracker.git
cd my-body-tracker
```

---

### 2. **Configurer les variables d'environnement**

Le fichier `.env.exemple` fournit un modèle pour les variables d'environnement nécessaires. Copie ce fichier et adapte-le selon tes besoins :

```bash
cp .env.exemple .env
```

Modifie le fichier `.env` pour y renseigner les valeurs appropriées (par exemple, les ports, les clés API, etc.).

---

### 3. **Utiliser Docker pour l'exécution**

Le projet inclut un `docker-compose.yml` et un `Dockerfile`, facilitant le déploiement via Docker. Assure-toi d'avoir Docker et Docker Compose installés sur ta machine.

Pour construire et lancer les conteneurs :

```bash
docker-compose up --build
```

Cette commande construira l'image Docker définie dans le `Dockerfile` et démarrera les services spécifiés dans `docker-compose.yml`.

---

### 4. **Utiliser Make pour simplifier les commandes**

Un `Makefile` est présent pour automatiser certaines tâches. Par exemple, pour construire le projet :

```bash
make build
```

Consulte le contenu du `Makefile` pour découvrir d'autres commandes utiles, telles que `make run` ou `make test`.

---

### 5. **Structure du projet**

Le projet est structuré de manière modulaire :

- `cmd/` : Contient le point d'entrée principal de l'application.
- `api/` : Gère les routes et les contrôleurs de l'API.
- `cli/` : Implémente l'interface en ligne de commande.
- `internal/calculator/` : Regroupe la logique métier, notamment les calculs liés à la nutrition.

---

### 6. **Exécution de l'application**

Après avoir configuré les variables d'environnement et lancé les conteneurs Docker, l'application devrait être accessible via l'interface CLI ou les endpoints API, selon l'implémentation.

## Fonctionnalités de l'application

### 1. Informations personnelles de santé

En tant qu'utilisateur de l'application, vous devez fournir certaines informations personnelles de santé, telles que votre poids, votre taille, votre âge et votre niveau d'activité physique.
Ces données sont indispensables pour calculer des indicateurs clés comme l'**IMC**, l'**IMG** et vos besoins caloriques journaliers.

Grâce à la création d'un profil personnel, vous pouvez suivre l'évolution de votre condition physique au fil du temps.
En enregistrant régulièrement vos mesures, l'application met à jour vos indicateurs de santé, vous permettant ainsi de visualiser vos progrès de manière claire et précise.

Conçue pour une utilisation locale, l'application garantit la confidentialité de vos informations : aucune donnée personnelle n'est transmise à un serveur.
Elle fonctionne de manière autonome sur votre machine, assurant ainsi une protection totale de votre vie privée.

#### 1.1. Calcul de l'IMC

L'**Indice de Masse Corporelle** (IMC) est un indicateur permettant d'évaluer la corpulence d'une personne en fonction de son poids et de sa taille. Il est calculé à l'aide de la formule suivante :

$$\text{IMC} = \frac{\text{Masse (Kg)}}{\text{(Taille (m))}^2}$$

Cet indicateur est utilisé pour identifier les risques liés à une insuffisance pondérale ou à un excès de poids.

#### 1.2. Calcul de l'IMG

L'**Indice de Masse Grasse** (IMG) **n'est pas un indicateur fiable à 100 %** s'il est utilisé seul. Il **ne s'applique pas** aux enfants, aux adolescents de moins de 15 ans, aux personnes de plus de 50 ans, aux personnes très musclées, ni aux femmes enceintes.
[_Source Allo Docteurs_](https://www.allodocteurs.fr/maladies/obesite/indice-de-masse-grasse-img_195.html)

$$\text{IMG (%)} = (1.20∗\text{IMC}) + (0.23∗\text{Age (an)}) − (10.8∗\text{Sexe [Homme : 1, Femme : 0]}) − 5.4$$

**Tableau de classification de l'IMG :**

| **Catégorie**            | **Femmes (IMG)**     | **Hommes (IMG)**     |
| ------------------------ | -------------------- | -------------------- |
| Masse grasse trop faible | Inférieur à <25 %    | Inférieur à <15 %    |
| Masse grasse normale     | Entre >25 % et <30 % | Entre >15 % et <20 % |
| Masse grasse trop élevée | Supérieur à >30 %    | Supérieur à >20 %    |

### 2. Suivi des dépenses énergétiques

#### 2.1. Calcul du métabolisme de base (MB)

Le suivi des dépenses énergétiques permet de calculer les besoins caloriques journaliers, en fonction du poids, de la taille, de l'âge, du sexe et du niveau d'activité physique.
Les formules utilisées prennent en compte les spécificités physiologiques des hommes et des femmes.
[_Source Nutri&CO_](https://nutriandco.com/fr/pages/calcul-apport-calorique-journalier)

- **Pour les hommes** $$\text{MB (kcal)} = [1.083 \times \text{Poids (kg)} \times 0.48 \times \text{Taille (m)} \times 0.50 \times \text{Âge (an)} - 0.13] \times \frac{1000}{4.1855}$$
- **Pour les femmes** $$\text{MB (kcal)} = [0.963 \times \text{Poids (kg)} \times 0.48 \times \text{Taille (m)} \times 0.50 \times \text{Âge (an)} - 0.13] \times \frac{1000}{4.1855}$$

#### 2.2. Calcul du niveau d'activité physique (NAP)

Le coefficient ou niveau d'activité physique (NAP), varie selon l'intensité et la fréquence des activités physiques pratiquées.

$$\text{NAP} = 1.15 + \text{Intensité (J)} \times 0.3946$$

**Tableau des coefficients NAP :**

| Intensité d\'activité physique                                | Coefficient NAP |
| ------------------------------------------------------------- | --------------- |
| Sédentaire (peu ou pas d\'activité physique)                  | 1.0 à 1.3       |
| Activité légère (1 à 3 jours par semaine)                     | 1.375           |
| Activité modérée (3 à 5 jours par semaine)                    | 1.55            |
| Activité intense (5 à 6 jours par semaine)                    | 1.725           |
| Activité très intense (entraînement quotidien ou biquotidien) | 1.9             |

### 3. Suivi de l'alimentation quotidienne

L'application permet à l'utilisateur d'ajouter les aliments consommés à chaque repas de la journée : _petit-déjeuner_, _déjeuner_, _dîner_, _collation_.
Elle récupère automatiquement leurs informations nutritionnelles, telles que les glucides, protéines, lipides.
Ces données sont obtenues grâce à la base de données [FoodData Central (FDC)](https://fdc.nal.usda.gov/), fournie par le [National Agricultural Library (NAL)](https://www.nal.usda.gov/) et le [United States Department of Agriculture (USDA)](https://www.usda.gov/).

Les informations nutritionnelles sont présentées sous forme de tableau, offrant une vue clair des apports de chaque aliment.

### 4. Objectif et suivi de progression

L'utilisateur peut définir un objectif nutritionnel personnalisé en lien avec le suivi de son alimentation quotidienne. L'application permet de suivre l'évolution en fonction des apports caloriques et des dépenses énergétiques, tout en offrant une visualisation détaillée des proportions de macronutriments (glucides, lipides, protéines).

Ces fonctionnalités aident l'utilisateur à mieux gérer son alimentation et à atteindre ou maintenir un poids.

#### 4.1. Calcul du Besoin Énergétique Journalier (BEJ)

Le besoin énergétique journalier (BEJ) est obtenu en multipliant les dépenses du métabolisme de base (MB) par un coefficient d'activité physique (NAP - Niveau d’Activité Physique).

$$\text{BEJ} = \text{MB} \times \text{NAP}$$

#### 4.2. Répartition des Macronutriments

Une fois le **BEJ** calculé, il est possible de répartir les macronutriments (glucides, lipides, protéines) selon l'objectif visé.

**Formules générales pour les apports en macronutriments** :

$$\text{Glucides (g)} = \left( \frac{\text{Pourcentage de glucides} \times BEJ}{4} \right)$$

$$\text{Protéines (g)} = \left( \frac{\text{Pourcentage de protéines} \times BEJ}{4} \right)$$

$$\text{Lipides (g)} = \left( \frac{\text{Pourcentage de lipides} \times BEJ}{9} \right)$$

**Apport calorique par gramme de macronutriments** :

- 1 g de **glucides** = **4 kcal**
- 1 g de **protéines** = **4 kcal**
- 1 g de **lipides** = **9 kcal**

**Ratios recommandés selon l'objectif**

| **Objectif**                  | **Glucides (%)** | **Lipides (%)** | **Protéines (%)** |
| ----------------------------- | ---------------- | --------------- | ----------------- |
| **Équilibre énergétique**     | 50               | 35              | 15                |
| **Perte de poids**            | 35               | 30              | 35                |
| **Régime cétogène**           | 10               | 50              | 40                |
| **Régime hyperprotéiné**      | 10               | 25              | 65                |
| **Prise de masse musculaire** | 50               | 10              | 40                |
| **Sportif de haut niveau**    | 50               | 20              | 30                |

#### 4.3. Calcul des Apports selon le Poids de Corps

Dans certains cas (prise de masse musculaire, diète sportive), les macronutriments sont calculés en **grammes par kg de poids corporel**.

Voici une liste des ratios recommandés pour la prise de masse musculaire :

- $`\text{Glucides (g)} = 4  \times \text{Poids (kg)}`$

- $`\text{Protéines (g)} = 2 \times \text{Poids (kg)}`$

- $`\text{Lipides (g)} = 1 \times \text{Poids (kg)}`$
