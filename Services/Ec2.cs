using System;
using System.Collections.Generic;
using System.Dynamic;
using System.IO;
using System.Linq;
using System.Reflection;
using System.Threading.Tasks;
using Amazon.EC2;
using Amazon.EC2.Model;
using CsvHelper;
using CsvLibrary;
using Dynamitey;

namespace Services
{
    public class Ec2
    {
        private readonly AmazonEC2Client _amazonEc2Client = new AmazonEC2Client();
        private Csv Csv { get; }

        public Ec2(Csv csv)
        {
            Csv = csv;
        }
        
        private async Task<List<object>> GenerateEc2List()
        {
            List<object> csvRecords = new List<object>();
            DescribeInstancesResponse instances = await _amazonEc2Client.DescribeInstancesAsync();
            var instancesFiltered = instances.Reservations.Select(a => a.Instances);

            foreach (var i in instancesFiltered)
            {
                foreach (var instance in i)
                {
                    csvRecords.Add(new {Id = instance.InstanceId, Name = instance.Tags.
                        FirstOrDefault(a => a.Key == "Name")
                        ?.Value});
                }
            }

            return csvRecords;
        }

        public void SaveToCsv()
        {
            Csv.Write("ec2.csv", GenerateEc2List().Result);
        }

        public void ReadCsv()
        {
            List<CreateTagsRequest> tagsRequests = new List<CreateTagsRequest>();
            List<Tag> tags = new List<Tag>();
            var csv = Csv.Read("ec2.csv");
            Dynamitey.GetMemberNames(csv[0]);
            foreach (var record in csv)
            {
                Console.WriteLine(record);
            }
            /*foreach (var filed in csv)
            {
                foreach (var t in csv)
                {
                    /*Console.WriteLine($"{t.Key} = {t.Value}");#1#
                    Console.WriteLine(t.Key == "Id" ? t.Value.ToString(): "i-99999");
                    tagsRequests.Add(new CreateTagsRequest
                    {
                        Resources = new List<string>
                        {
                            t.Key == "Id" ? t.Value.ToString(): null
                        },
                        Tags = new List<Tag>
                        {
                        }
                    });
                }
            }*/
        }
    }
}