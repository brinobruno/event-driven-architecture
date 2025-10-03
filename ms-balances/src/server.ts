import Fastify from 'fastify'
import pg from '@fastify/postgres'
import { kafka } from './external/kafka.js'

const fastify = Fastify({
  logger: true
})

fastify.register(pg, {
  connectionString: 'postgres://postgres:postgres@postgres:5432/wallet'
})

// Not using migration for simplicity on this project
fastify.ready(async () => {
  await fastify.pg.query(`
    CREATE TABLE IF NOT EXISTS balances (
      id SERIAL PRIMARY KEY,
      account_id VARCHAR(255) NOT NULL,
      balance DECIMAL(10, 2) NOT NULL
    )
  `)
})

fastify.get('/', async function handler (request, reply) {
  return { hello: 'world' }
})

const consumer = kafka.consumer({ groupId: 'walletcore' })
await consumer.connect()
await consumer.subscribe({ topic: 'balances', fromBeginning: true })
await consumer.run({
  eachMessage: async ({ topic, partition, message }) => {
    console.log({
      topic,
      partition,
      offset: message.offset,
      value: message.value?.toString(),
    })
  },
})

try {
  await fastify.listen({ port: 3000 })
} catch (err) {
  fastify.log.error(err)
  process.exit(1)
}