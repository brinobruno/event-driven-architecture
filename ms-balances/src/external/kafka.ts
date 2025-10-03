import { Kafka } from 'kafkajs'

export const kafka = new Kafka({
  clientId: 'walletcore',
  brokers: ['kafka:29092'],
})