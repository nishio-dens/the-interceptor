# need https://github.com/nishio-dens/convergence to setup your tables

create_table :interceptor_buckets do |t|
  t.int :id, primary_key: true, extra: :auto_increment
  t.varchar :name
  t.int :read_bucket_id
  t.int :write_bucket_id

  t.datetime :created_at, null: true
  t.datetime :updated_at, null: true
  t.datetime :deleted_at, null: true

  t.index :name, unique: true
  t.index :read_bucket_id
  t.index :write_bucket_id
end

create_table :s3_buckets do |t|
  t.int :id, primary_key: true, extra: :auto_increment

  t.varchar :bucket_name
  t.varchar :bucket_access_key, null: true
  t.varchar :bucket_access_secret, null: true
  t.varchar :bucket_region, default: "us-east-1"
  t.boolean :bucket_disable_ssl, default: true
  t.varchar :bucket_endpoint, null: true

  t.datetime :created_at, null: true
  t.datetime :updated_at, null: true
  t.datetime :deleted_at, null: true
end

create_table :interceptor_objects, comment: "Intercept Object Histories" do |t|
  t.int :id, primary_key: true, extra: :auto_increment
  t.int :virtual_bucket_id
  t.int :real_bucket_id

  t.varchar :key, limit: 4096
  t.int :size, default: 0

  t.datetime :created_at, null: true
  t.datetime :updated_at, null: true
  t.datetime :deleted_at, null: true

  t.index :key, length: 255
end

create_table :users do |t|
  t.int :id, primary_key: true, extra: :auto_increment
  t.varchar :username
  t.varchar :access_key, null: true
  t.varchar :access_secret, null: true

  t.datetime :created_at, null: true
  t.datetime :updated_at, null: true
  t.datetime :deleted_at, null: true
end
