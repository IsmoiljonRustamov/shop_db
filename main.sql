CREATE TABLE "categories"(
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "image_url" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE "products"(
    "id" SERIAL PRIMARY KEY,
    "category_id" INTEGER NOT NULL REFERENCES categories(id) on delete restrict ,
    "name" VARCHAR(255) NOT NULL,
    "price" DECIMAL(18, 2) NOT NULL,
    "image_url" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE "customers"(
    "id" SERIAL PRIMARY KEY,
    "first_name" VARCHAR(255) NOT NULL,
    "last_name" VARCHAR(255) NOT NULL,
    "phone_number" VARCHAR(255) NOT NULL,
    "gender" INTEGER NOT NULL,
    "birth_date" DATE NOT NULL,
    "created_at" TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
    "deleted_at" TIMESTAMP
);
CREATE TABLE "orders"(
    "id" SERIAL PRIMARY KEY,
    "customer_id" INTEGER NOT NULL REFERENCES customers(id),
    "total_amount" DECIMAL(18, 2) NOT NULL,
    "created_at" TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "address" TEXT NOT NULL
);
CREATE TABLE "order_items"(
    "id" SERIAL PRIMARY KEY,
    "order_id" INTEGER NOT NULL REFERENCES orders(id),
    "product_id" INTEGER NOT NULL REFERENCES products(id),
    "count" INTEGER NOT NULL,
    "total_price" DECIMAL(18, 2) NOT NULL,
    "product_name" VARCHAR(255) NOT NULL
);
CREATE TABLE "product_images"(
    "id" SERIAL PRIMARY KEY,
    "image_url" INTEGER NOT NULL,
    "sequence_number" INTEGER NOT NULL,
    "product_id" INTEGER NOT NULL REFERENCES products(id)
);











CREATE TABLE "cars_db"(
    "id" SERIAL PRIMARY KEY,
    "car_name" VARCHAR(255) NOT NULL,
    "color" VARCHAR(255) NOT NULL,
    "price" DECIMAL(18, 2) NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
    "image_url" VARCHAR(255) NOT NULL
);

CREATE TABLE "cars_image"(
    "id" INTEGER NOT NULL,
    "image_url" INTEGER NOT NULL,
    "sequence_number" INTEGER NOT NULL
    "car_id" INTEGER NOT NULL REFERENCES cars_db(id)
);
